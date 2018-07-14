package midgard

import (
	"errors"
)

var ErrInvaildAggregateID = errors.New("aggregate id is empty")
var ErrNilAggregate = errors.New("aggregate is nil")

type EventRepository interface {
	Load(aggregate Aggregate, aggregateID string) error
	Save(aggregateID string, events ...Event) error
	SaveAndCommit(aggregateID string, events ...Event) error
	TxBegin()
	Commit() error
	Close()
}

type EventPublisher interface {
	Publisher
}

type Repository struct {
	store     EventStore
	publisher EventPublisher
}

func NewRepo(store EventStore, publisher EventPublisher) *Repository {
	return &Repository{
		store:     store,
		publisher: publisher,
	}
}

//Load aggregate by id and replay saved event
func (r *Repository) Load(aggregate Aggregate, aggregateID string) error {

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

// save events to aggregate and publish event
// todo 보상 event 추가
func (r *Repository) Save(aggregateID string, events ...Event) error {

	if err := checkEvents(events...); err != nil {
		return err
	}

	return r.store.Save(aggregateID, events...)
}

func (r *Repository) TxBegin() {
	r.store.TxBegin()
}

//todo events publish all
func (r *Repository) Commit() error {

	if err := r.store.Commit(); err != nil {
		return err
	}
	//
	//if r.publisher != nil {
	//
	//	for _, event := range events {
	//		//Todo type implicit problem
	//		err := r.publisher.Publish("Event", event.GetType(), event)
	//		if err != nil {
	//			return errors.New("need roll back")
	//		}
	//	}
	//}

	return nil
}

func (r *Repository) SaveAndCommit(aggregateID string, events ...Event) error {

	if err := checkEvents(events...); err != nil {
		return err
	}

	err := r.store.SaveAndCommit(aggregateID, events...)

	if err != nil {
		return err
	}

	if r.publisher != nil {

		for _, event := range events {
			//Todo type implicit problem
			err := r.publisher.Publish("Event", event.GetType(), event)
			if err != nil {
				return errors.New("need roll back")
			}
		}
	}

	return nil
}

func checkEvents(events ...Event) error {

	if len(events) == 0 {
		return errors.New("no events to save")
	}

	for _, event := range events {
		if event.GetType() == "" {
			return errors.New("all event need type")
		}
	}

	return nil
}

func (r *Repository) Close() {
	r.store.Close()
}
