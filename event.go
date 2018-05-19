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
type EventD struct {
	AggregateID string
	Version     int
	Time        time.Time
	Type        string
}

func (e EventD) GetType() string {
	return e.Type
}

func (e EventD) GetAggregateID() string {
	return e.AggregateID
}

func (e EventD) GetVersion() int {
	return e.Version
}
