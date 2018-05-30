package midgard

import "time"

//Event interface
type Event interface {
	Entity
	GetType() string
	GetVersion() int
}

//Providing default an Event implementation all event must have this EventModel
type EventModel struct {

	// ID of aggregate Root
	ID string

	// Specifies the version of the event. Modify the version if the structure of the event changes.
	Version int

	// Specifies topic when published through publisher
	Type string
	Time time.Time
}

func (e EventModel) GetId() string {
	return e.AggregateID
}

func (e EventModel) GetType() string {
	return e.Type
}

func (e EventModel) GetID() string {
	return e.ID
}

func (e EventModel) GetVersion() int {
	return e.Version
}
