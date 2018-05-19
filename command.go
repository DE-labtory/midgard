package eventsource

type Command struct {
}

//Command를 받아 Event를 발생시키는 Handler
type CommandHandler interface {
	Handle(command Command)
}
