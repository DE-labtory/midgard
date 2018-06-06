package mongodb

import (
	"gopkg.in/mgo.v2"
	"github.com/it-chain/midgard"
	"sync"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"encoding/json"
	"strings"
	"errors"
)

var ErrNilEvents = errors.New("no event history exist")

type SerializedEvent struct {
	Type string
	Data []byte
}

type History []SerializedEvent

type Document struct {
	AggregateID string 			`bson:"aggregate_id"`
	History 					`bson:"history"`
}

func (d *Document) appendEvent(serializedEvent SerializedEvent) {
	d.History = append(d.History, serializedEvent)
}

func (d *Document) getHistory() History {
	return d.History
}

type Store struct {
	name string
	mux *sync.RWMutex
	*mgo.Session
	mgo.Index
	serializer EventSerializer
}

func NewEventStore(path string, db string, serializer EventSerializer) midgard.EventStore {
	s, err := mgo.Dial(path)

	if err != nil {
		return nil
	}
	return &Store{
		name: db,
		mux: &sync.RWMutex{},
		Session: s,
		serializer: serializer,
		Index: mgo.Index{
			Key:        []string{"aggregate_id"},
			Unique:     true, 		// Prevent two documents from having the same index key
			// DropDups:   false, 	// Drop documents with the same index key as a previously indexed one
			Background: true, 		// Build index in background and return immediately
			Sparse:     true, 		// Only index documents containing the Key fields
		},
	}

}

//Save Events to mongodb
func (s Store) Save(aggregateID string, events ...midgard.Event) error {
	s.mux.Lock()
	session := s.getFreshSession()

	defer func() {
		s.mux.Unlock()
		session.Close()
	}()

	document, err := s.getDocument(aggregateID)

	if err != nil {
		document = &Document{
			AggregateID: aggregateID,
			History: []SerializedEvent{},
		}
	}

	for _, event := range events {
		serializedEvent, err := s.serializer.Marshal(event)
		if err != nil {
			return err
		}
		document.appendEvent(serializedEvent)
	}

	c := session.DB(s.name).C("events")
	c.EnsureIndex(s.Index)

	_, err = c.Upsert(bson.M{"aggregate_id": aggregateID}, document)
	return err
}

//Load Aggregate Event from leveldb
func (s Store) Load(aggregateID string) ([]midgard.Event, error) {
	s.mux.Lock()
	session := s.getFreshSession()

	defer func() {
		s.mux.Unlock()
		session.Close()
	}()

	document, err := s.getDocument(aggregateID)

	if err != nil {
		return nil, err
	}

	events := make([]midgard.Event, 0)

	for _, v := range document.getHistory() {
		event, err := s.serializer.Unmarshal(v)

		if err != nil {
			return []midgard.Event{}, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (s Store) getDocument(aggregateID string) (*Document, error) {
	var document = Document{}

	c := s.Session.DB(s.name).C("events")
	err := c.Find(bson.M{"aggregate_id": aggregateID}).One(&document)

	return &document, err
}

// When do multithreaded work, we want to be thread-safe
// open another session from the database pool
func (s Store) getFreshSession() *mgo.Session {
	return s.Session.Copy()
}

type EventSerializer interface {
	// MarshalEvent converts an Event to a Record
	Marshal(event midgard.Event) (SerializedEvent, error)

	// UnmarshalEvent converts an Event backed into a Record
	Unmarshal(serializedEvent SerializedEvent) (midgard.Event, error)
}

type JSONSerializer struct {
	eventTypes map[string]reflect.Type
}

func NewSerializer(events ...midgard.Event) EventSerializer {

	s := &JSONSerializer{
		eventTypes: make(map[string]reflect.Type),
	}

	s.Register(events...)

	return s
}

func (j *JSONSerializer) Register(events ...midgard.Event) {

	for _, event := range events {
		rawType, name := GetTypeName(event)
		j.eventTypes[name] = rawType
	}
}

func (j *JSONSerializer) Marshal(e midgard.Event) (SerializedEvent, error) {

	serializedEvent := SerializedEvent{}
	_, name := GetTypeName(e)
	serializedEvent.Type = name

	data, err := json.Marshal(e)

	if err != nil {
		return SerializedEvent{}, err
	}

	serializedEvent.Data = data

	return serializedEvent, nil
}

func (j *JSONSerializer) Unmarshal(serializedEvent SerializedEvent) (midgard.Event, error) {

	t, ok := j.eventTypes[serializedEvent.Type]

	if !ok {
		return nil, errors.New(fmt.Sprintf("unbound event type, %v", serializedEvent.Type))
	}

	v := reflect.New(t).Interface()

	err := json.Unmarshal(serializedEvent.Data, v)
	if err != nil {
		return nil, err
	}

	return v.(midgard.Event), nil
}

func GetTypeName(source interface{}) (reflect.Type, string) {

	rawType := reflect.TypeOf(source)

	if rawType.Kind() == reflect.Ptr {
		rawType = rawType.Elem()
	}

	name := rawType.String()
	parts := strings.Split(name, ".")
	return rawType, parts[1]
}