package store

import "github.com/it-chain/eventsource"

type EventStore interface {
	Save(aggregateID string, events ...eventsource.Event) error
	Load(aggregateID string) ([]eventsource.Event, error)
}
