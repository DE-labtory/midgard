package leveldb

import (
	"encoding/json"
	"errors"

	"sync"

	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/store"
)

var ErrNilEvents = errors.New("no event history exist")
var ErrGetValue = errors.New("fail to get value from leveldb")

type History []store.SerializedEvent

//Leveldb store implementing store interface
type Store struct {
	mux        *sync.RWMutex
	db         *leveldbwrapper.DB
	serializer store.EventSerializer
}

func NewEventStore(path string, serializer store.EventSerializer) midgard.EventStore {

	db := leveldbwrapper.CreateNewDB(path)
	db.Open()

	return &Store{
		db:         db,
		mux:        &sync.RWMutex{},
		serializer: serializer,
	}
}

//Save Events to leveldb(key is aggregateID)
func (s Store) Save(aggregateID string, events ...midgard.Event) error {

	s.mux.Lock()
	defer s.mux.Unlock()

	history, err := s.getHistory(aggregateID)

	if err != nil {
		return err
	}

	//new history
	if history == nil {
		history = &History{}
	}

	for _, event := range events {
		serializedEvent, err := s.serializer.Marshal(event)

		if err != nil {
			return err
		}

		*history = append(*history, serializedEvent)
	}

	historyValue, err := json.Marshal(history)

	if err != nil {
		return err
	}

	return s.db.Put([]byte(aggregateID), historyValue, true)
}

//Load Aggregate Event from leveldb
func (s Store) Load(aggregateID string) ([]midgard.Event, error) {

	history, err := s.getHistory(aggregateID)

	if err != nil {
		return nil, err
	}

	//new history
	if history == nil {
		return nil, ErrNilEvents
	}

	events := make([]midgard.Event, 0)

	for _, value := range *history {
		event, err := s.serializer.Unmarshal(value)

		if err != nil {
			return []midgard.Event{}, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (s Store) getHistory(aggregateID string) (*History, error) {

	var history = &History{}

	historyValue, err := s.db.Get([]byte(aggregateID))

	if err != nil {
		return nil, ErrGetValue
	}

	//history does not exist
	if historyValue == nil {
		return nil, nil
	}

	err = json.Unmarshal(historyValue, history)

	if err != nil {
		return nil, err
	}

	return history, nil
}

func (s Store) Close() {
	s.db.Close()
}
