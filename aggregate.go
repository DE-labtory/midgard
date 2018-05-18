package eventsource

//Domain Entitiy의 Aggregate
type Aggregate interface{
	CommandHandler
	EventReceptor
}

//Command를 받아 Event를 발생시키는 Handler
type CommandHandler interface{
	Handle(command Command)
}

//Event를 받아 상태를 변화시키는 Receptor
type EventReceptor interface{
	On(event Event)
}