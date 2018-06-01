package midgard

//Entity has unique ID
type Entity interface {
	GetID() string
}

//Domain Entitiy의 Aggregate
type Aggregate interface {
	EventReceptor
	Entity
}

type AggregateModel struct {
	// ID of aggregate Root
	ID string
}

func (aggregate AggregateModel) GetID() string {
	return aggregate.ID
}

//Event를 받아 상태를 변화시키는 Receptor
type EventReceptor interface {
	On(event Event) error
}
