package mongodb

import (
	"sync"

	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

type History []store.SerializedEvent

type Document struct {
	AggregateID string `bson:"aggregate_id"`
	History     `bson:"history"`
}

func (d *Document) appendEvent(serializedEvent store.SerializedEvent) {
	d.History = append(d.History, serializedEvent)
}

func (d *Document) getHistory() History {
	return d.History
}

type Store struct {
	name string
	mux  *sync.RWMutex
	*mgo.Session
	mgo.Index
	serializer store.EventSerializer
	txEvents   map[string][]midgard.Event
}

func NewEventStore(url string, db string, serializer store.EventSerializer) (*Store, error) {

	s, err := mgo.Dial(url)

	if err != nil {
		return nil, err
	}

	return &Store{
		name:       db,
		mux:        &sync.RWMutex{},
		Session:    s,
		serializer: serializer,
		txEvents:   make(map[string][]midgard.Event, 0),
		Index: mgo.Index{
			Key:    []string{"aggregate_id"},
			Unique: true, // Prevent two documents from having the same index key
			// DropDups:   false, 	// Drop documents with the same index key as a previously indexed one
			Background: true, // Build index in background and return immediately
			Sparse:     true, // Only index documents containing the Key fields
		},
	}, nil
}

func (s Store) SaveAndCommit(aggregateID string, events ...midgard.Event) error {
	return s.save(aggregateID, events...)
}

func (s *Store) TxBegin() {

	s.mux.Lock()

	for key, _ := range s.txEvents {
		delete(s.txEvents, key)
	}
}

func (s *Store) Commit() error {

	session := s.getFreshSession()

	defer func() {
		s.mux.Unlock()
	}()

	txns := make([]txn.Op, 0)

	for id, events := range s.txEvents {

		txnOp, err := s.createTxnOp(id, events...)

		if err != nil {
			return err
		}

		txns = append(txns, txnOp)
	}

	c := session.DB(s.name).C("events")
	c.EnsureIndex(s.Index)

	runner := txn.NewRunner(c)

	return runner.Run(txns, "", nil)
}

func (s *Store) createTxnOp(aggregateID string, events ...midgard.Event) (txn.Op, error) {

	session := s.getFreshSession()
	defer session.Close()

	document, err := s.getDocument(aggregateID)

	if err == mgo.ErrNotFound {
		return s.insert(session, aggregateID, events...)
	}

	if err != nil {
		return txn.Op{}, err
	}

	return s.update(document, session, aggregateID, events...)
}

func (s *Store) insert(session *mgo.Session, aggregateID string, events ...midgard.Event) (txn.Op, error) {

	document := &Document{
		AggregateID: aggregateID,
		History:     []store.SerializedEvent{},
	}

	for _, event := range events {
		serializedEvent, err := s.serializer.Marshal(event)
		if err != nil {
			return txn.Op{}, err
		}
		document.appendEvent(serializedEvent)
	}

	return txn.Op{
		C:      "events",
		Id:     aggregateID,
		Assert: txn.DocMissing,
		Insert: document,
	}, nil
}

func (s *Store) update(document *Document, session *mgo.Session, aggregateID string, events ...midgard.Event) (txn.Op, error) {

	for _, event := range events {
		serializedEvent, err := s.serializer.Marshal(event)
		if err != nil {
			return txn.Op{}, err
		}
		document.appendEvent(serializedEvent)
	}

	return txn.Op{
		C:      "events",
		Id:     aggregateID,
		Assert: txn.DocExists,
		Update: bson.M{"$set": document},
	}, nil
}

func (s *Store) Save(aggregateID string, events ...midgard.Event) error {

	storedEvents, ok := s.txEvents[aggregateID]

	if ok {
		storedEvents = append(storedEvents, events...)
		s.txEvents[aggregateID] = storedEvents
	}

	s.txEvents[aggregateID] = events

	return nil
}

//Save Events to mongodb
func (s Store) save(aggregateID string, events ...midgard.Event) error {

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
			History:     []store.SerializedEvent{},
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

func (s Store) Close() {
	s.Session.Close()
}
