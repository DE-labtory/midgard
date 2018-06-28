package event_store

import (
	"sync"

	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/store"
	"github.com/it-chain/midgard/store/mongodb"
)

var once sync.Once

var instance *EventStore

type EventStore struct {
	repository midgard.EventRepository
	serializer store.EventSerializer
}

func InitMongoEventStore(url string, dbname string, publisher midgard.EventPublisher, events ...midgard.Event) *EventStore {

	once.Do(func() {

		instance = &EventStore{}
		instance.serializer = store.NewSerializer(events...)

		store, err := mongodb.NewEventStore(url, dbname, instance.serializer)

		if err != nil {
			panic(err)
		}

		r := midgard.NewRepo(store, publisher)

		instance.repository = r
	})

	// it-chain의 설정내용을 반환한다.
	return instance
}

func RegisterEvents(event ...midgard.Event) {
	instance.serializer.Register(event...)
}

func Save(aggregateID string, events midgard.Event) error {
	return instance.repository.Save(aggregateID, events)
}

func Load(aggregate midgard.Aggregate, aggregateID string) error {
	return instance.repository.Load(aggregate, aggregateID)
}

func Close() {
	instance.repository.Close()
}
