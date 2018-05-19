package store

import "github.com/it-chain/eventsource"

type EventStore interface{
	Save(aggregateID string, events ...eventsource.Event)
	Load(aggregateID string, aggregate eventsource.Aggregate) error
}