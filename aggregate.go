package eventsource

//Domain Entitiy의 Aggregate
type Aggregate interface {
	EventReceptor
	GetAggregateID() string
}

type AggregateModel struct {
	// ID of aggregate Root
	AggregateID string
}

func (aggregate AggregateModel) GetAggregateID() string {
	return aggregate.AggregateID
}

//Event를 받아 상태를 변화시키는 Receptor
type EventReceptor interface {
	On(event Event) error
}
