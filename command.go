package midgard

type Command interface {
	Entity
}

type CommandModel struct {
	// ID contains the aggregate id
	ID string
}

func (c CommandModel) GetID() string {
	return c.ID
}
