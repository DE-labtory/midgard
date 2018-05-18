package eventsource

import "time"

//Providing default an Event
type Event struct{
	AggregateID string
	Version int
	Time time.Time
	Type string
}