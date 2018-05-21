package eventsource

type Command interface {
	GetAggregateID() string
}

type CommandModel struct {
	// ID contains the aggregate id
	AggregateID string
}

func (c CommandModel) GetAggregateID() string {
	return c.AggregateID
}
