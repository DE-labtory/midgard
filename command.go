package eventsource

type Command interface {
	AggregateID() string
}

type CommandModel struct {
	// ID contains the aggregate id
	ID string
}

func (c CommandModel) AggregateID() string {
	return c.ID
}

//Command를 받아 Event를 발생시키는 Handler
type CommandHandler interface {
	Handle(command Command)
}
