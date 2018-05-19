package eventsource

//Domain Entitiy의 Aggregate
type Aggregate interface {
	EventReceptor
}

//Event를 받아 상태를 변화시키는 Receptor
type EventReceptor interface {
	On(event Event) error
}
