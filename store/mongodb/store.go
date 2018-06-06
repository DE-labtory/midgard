package mongodb

import (
	"gopkg.in/mgo.v2"
	"github.com/it-chain/midgard"
	"sync"
	"gopkg.in/mgo.v2/bson"
	"github.com/it-chain/midgard/store"
)


type History []store.SerializedEvent

type Document struct {
	AggregateID string 			`bson:"aggregate_id"`
	History 					`bson:"history"`
}

func (d *Document) appendEvent(serializedEvent store.SerializedEvent) {
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
	serializer store.EventSerializer
}

func NewEventStore(path string, db string, serializer store.EventSerializer) midgard.EventStore {
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
			History: []store.SerializedEvent{},
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

