package main

import (
	"errors"
	"fmt"
	"os"

	"sync"

	"log"

	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/it-chain/midgard/store/leveldb"
)

// aggregate
type User struct {
	Name string
	midgard.AggregateModel
}

func (u User) On(event midgard.Event) error {

	switch v := event.(type) {

	case *UserCreatedEvent:
		u.AggregateID = v.AggregateID

	case *UserNameUpdatedEvent:
		u.Name = v.Name

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

// Command
type UserCreateCommand struct {
	midgard.CommandModel
}

type UserNameUpdateCommand struct {
	midgard.CommandModel
	Name string
}

// Event
type UserCreatedEvent struct {
	midgard.EventModel
}

type UserNameUpdatedEvent struct {
	midgard.EventModel
	Name string
}

// CommandHandler
type UserCommandHandler struct {
	eventRepository *midgard.Repository
}

func (u UserCommandHandler) UserCreate(command UserCreateCommand) {

	log.Printf("received UserCreateCommand [%s]", command)
	events := make([]midgard.Event, 0)
	events = append(events, UserCreatedEvent{
		midgard.EventModel{
			AggregateID: "123",
			Type:        "User",
		},
	})

	err := u.eventRepository.Save(command.GetAggregateID(), events...)
	if err != nil {
		panic(err)
	}
}

func (u UserCommandHandler) UserNameUpdate(command UserNameUpdateCommand) {

	log.Printf("received UserUpdateCommand [%s]", command)

	user := User{}

	u.eventRepository.Load(user, command.AggregateID)

	events := make([]midgard.Event, 0)
	events = append(events, UserNameUpdatedEvent{
		midgard.EventModel{
			AggregateID: "123",
			Type:        "User",
		}, "Jun",
	})

	err := u.eventRepository.Save(command.GetAggregateID(), events...)

	if err != nil {
		panic(err)
	}
}

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	path := "test"

	c := rabbitmq.Connect("")
	store := leveldb.NewEventStore(path, leveldb.NewSerializer(UserCreatedEvent{}))
	repo := midgard.NewRepo(store, c)

	defer os.RemoveAll(path)

	err := c.Subscribe("Command", "User", &UserCommandHandler{eventRepository: repo})
	if err != nil {
		panic(err)
	}

	wg.Wait()
}
