package midgard

type EventStore interface {
	Save(aggregateID string, events ...Event) error
	Load(aggregateID string) ([]Event, error)
	SaveAndCommit(aggregateID string, events ...Event) error
	TxBegin()
	Commit() error
	Close()
}
