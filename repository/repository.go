package repository

import (
	"errors"

	"github.com/it-chain/eventsource"
	"github.com/it-chain/eventsource/store"
)

var ErrInvaildAggregateID = errors.New("aggregate id is empty")
var ErrNilAggregate = errors.New("aggregate is nil")

type Repository struct {
	store store.EventStore
}

func New(store store.EventStore) *Repository {
	return &Repository{
		store: store,
	}
}

func (r *Repository) Load(aggregate eventsource.Aggregate, aggregateID string) error {

	if aggregateID == "" {
		return ErrInvaildAggregateID
	}

	if aggregate == nil {
		return ErrNilAggregate
	}

	events, err := r.store.Load(aggregateID)

	if err != nil {
		return err
	}

	for _, event := range events {
		err = aggregate.On(event)

		if err != nil {
			return errors.New("fail to ")
		}
	}

	return nil
}

func (r *Repository) Save(aggregateID string, events ...eventsource.Event) error {

	if len(events) == 0 {
		return errors.New("no events to save")
	}

	return r.store.Save(aggregateID, events...)
}
