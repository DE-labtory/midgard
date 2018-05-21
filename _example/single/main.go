package main

import (
	"os"

	"fmt"

	"sync"

	"github.com/it-chain/eventsource"
	"github.com/it-chain/eventsource/bus/rabbitmq"
	"github.com/it-chain/eventsource/store/leveldb"
)

var wg = sync.WaitGroup{}

// aggregate
type User struct {
	name string
	eventsource.AggregateModel
}

// Command
type UserCreateCommand struct {
	eventsource.CommandModel
}

// Event
type UserCreatedEvent struct {
	eventsource.EventModel
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
	eventRepository *eventsource.Repository
}

func (u UserCommandHandler) UserCreated(command UserCreateCommand) {

	events := make([]eventsource.Event, 0)
	events = append(events, UserCreatedEvent{
		eventsource.EventModel{
			AggregateID: "123",
			Type:        "User",
		},
	})

	err := u.eventRepository.Save(command.GetAggregateID(), events...)
	if err != nil {
		panic(err)
	}
	wg.Done()
}

func main() {

	path := "test"

	c := rabbitmq.Connect("")
	store := leveldb.NewEventStore(path, leveldb.NewSerializer(UserCreatedEvent{}))
	r := eventsource.NewRepo(store, c)

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
