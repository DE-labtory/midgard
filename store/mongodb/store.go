package mongodb

import (
	"gopkg.in/mgo.v2"
	"github.com/it-chain/midgard"
	"sync"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type SerializedEvent struct {
	Type string
	Data []byte
}

type History  struct {
	AggregateID string 			`bson:"aggrgate_id"`
	Events []SerializedEvent 	`bson:"events"`
}

type Store struct {
	name string
	mux *sync.RWMutex
	*mgo.Session
	mgo.Index
}

func NewEventStore(path string, name string) midgard.EventStore {
	s, err := mgo.Dial(path)

	if err != nil {
		return nil
	}
	return &Store{
		name: name,
		mux: &sync.RWMutex{},
		Session: s,
		Index: mgo.Index{
			Key:        []string{"ID"},
			Unique:     false, 		// Prevent two documents from having the same index key
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

	c := session.DB(s.name).C("events");
	err := c.EnsureIndex(s.Index)
	if err != nil {
		return err
	}

	for _, event := range events {
		fmt.Println(event)
		err := c.Insert(event)

		if err != nil {
			return err
		}
	}

	return nil
}

//Load Aggregate Event from leveldb
func (s Store) Load(aggregateID string) ([]midgard.Event, error) {
	s.mux.Lock()
	session := s.getFreshSession()

	defer func() {
		s.mux.Unlock()
		session.Close()
	}()

	return nil, nil
}

func (s Store) getHistory(aggregateID string) (*History, error) {
	var history = History{}

	c := s.Session.DB(s.name).C("events")
	err := c.Find(bson.M{"aggregate_id": aggregateID}).One(&history)

	return &history, err
}

// When do multithreaded work, we want to be thread-safe
// open another session from the database pool
func (s Store) getFreshSession() *mgo.Session {
	return s.Session.Copy()
}