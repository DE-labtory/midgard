package mongodb

import (
	"gopkg.in/mgo.v2"
	"github.com/it-chain/midgard"
	"log"
	"sync"
)

type Store struct {
	mux *sync.RWMutex
	*mgo.Session
}

func NewEventStore(path string) midgard.EventStore {
	s, err := mgo.Dial(path)

	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &Store{
		Session: s,
		mux: &sync.RWMutex{},
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

	c := session.DB("midgard").C("events");

	for _, event := range events {
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

// When we do multithreaded work, we want to be thread-safe
// open another session from the database pool
func (s Store) getFreshSession() *mgo.Session {
	return s.Session.Copy()
}