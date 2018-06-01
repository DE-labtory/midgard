package main

import (
	"os"

	"fmt"

	"sync"

	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/it-chain/midgard/store/leveldb"
)

var wg = sync.WaitGroup{}

// aggregate
type User struct {
	name string
	midgard.AggregateModel
}

// Command
type UserCreateCommand struct {
	midgard.CommandModel
}

// Event
type UserCreatedEvent struct {
	midgard.EventModel
}

// EventHandler
type UserEventHandler struct {
}

func (u UserEventHandler) UserCreate(event UserCreatedEvent) {
	fmt.Println(event)
	wg.Done()
}

// CommandHandler
type UserCommandHandler struct {
	eventRepository *midgard.Repository
}

func (u UserCommandHandler) UserCreated(command UserCreateCommand) {

	events := make([]midgard.Event, 0)
	events = append(events, UserCreatedEvent{
		midgard.EventModel{
			ID:   "123",
			Type: "User",
		},
	})

	err := u.eventRepository.Save(command.GetID(), events...)
	if err != nil {
		panic(err)
	}
	wg.Done()
}

func main() {

	path := "test"

	c := rabbitmq.Connect("")
	store := leveldb.NewEventStore(path, leveldb.NewSerializer(UserCreatedEvent{}))
	r := midgard.NewRepo(store, c)

	defer os.RemoveAll(path)

	err := c.Subscribe("event", "User", &UserEventHandler{})
	err = c.Subscribe("command", "User", &UserCommandHandler{eventRepository: r})

	if err != nil {
		panic(err)
	}

	wg.Add(2)

	err = c.Publish("command", "User", UserCreateCommand{})

	if err != nil {
		panic(err)
	}

	wg.Wait()
}
