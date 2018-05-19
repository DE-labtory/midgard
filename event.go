package eventsource

import (
	"time"
)

type Event interface {
	GetType() string
	GetAggregateID() string
	GetVersion() int
}

//Providing default an Event implementation
type EventModel struct {
	AggregateID string
	Version     int
	Time        time.Time
	Type        string
}

func (e EventModel) GetType() string {
	return e.Type
}

func (e EventModel) GetAggregateID() string {
	return e.AggregateID
}

func (e EventModel) GetVersion() int {
	return e.Version
}
