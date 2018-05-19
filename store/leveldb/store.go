package leveldb

import (
	"encoding/json"
	"errors"
	"sync"

	"reflect"

	"github.com/it-chain/eventsource"
	"github.com/it-chain/eventsource/store"
	"github.com/it-chain/leveldb-wrapper"
)

var ErrNilEvents = errors.New("no event history exist")
var ErrGetValue = errors.New("fail to get value from leveldb")

type SerializedEvent []byte
type History []SerializedEvent

//Leveldb store implementing store interface
type Store struct {
	mux        *sync.RWMutex
	db         *leveldbwrapper.DB
	serializer EventSerializer
}

func NewEventStore(path string, serializer EventSerializer) store.EventStore {

	db := leveldbwrapper.CreateNewDB(path)
	db.Open()

	return &Store{
		db:         db,
		mux:        &sync.RWMutex{},
		serializer: serializer,
	}
}

//Save Events to leveldb(key is aggregateID)
func (s Store) Save(aggregateID string, events ...eventsource.Event) error {

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

	*history = append(*history, events...)
	historyValue, err := json.Marshal(history)

	if err != nil {
		return err
	}

	return s.db.Put([]byte(aggregateID), historyValue, true)
}

//Load Aggregate Event from leveldb
func (s Store) Load(aggregateID string) ([]eventsource.Event, error) {

	history, err := s.getHistory(aggregateID)

	if err != nil {
		return nil, err
	}

	//new history
	if history == nil {
		return nil, ErrNilEvents
	}

	return *history, nil
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

type EventSerializer interface {
	// MarshalEvent converts an Event to a Record
	Marshal(event eventsource.Event) (SerializedEvent, error)

	// UnmarshalEvent converts an Event backed into a Record
	Unmarshal(serializedEvent SerializedEvent) (eventsource.Event, error)
}

type JSONSerializer struct {
	eventTypes map[string]reflect.Type
}
