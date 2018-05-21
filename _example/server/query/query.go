package main

import (
	"fmt"

	"sync"

	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
)

type UserCreatedEvent struct {
	midgard.EventModel
}

type UserNameUpdatedEvent struct {
	midgard.EventModel
	Name string
}

type UserEventHandler struct {
}

func (u UserEventHandler) HandleCreatedEvent(event UserCreatedEvent) {
	fmt.Println(event)
}

func (u UserEventHandler) HandleNameUpdatedEvent(event UserNameUpdatedEvent) {
	fmt.Println(event)
}

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)

	c := rabbitmq.Connect("")
	err := c.Subscribe("Event", "User", &UserEventHandler{})

	if err != nil {
		panic(err.Error())
	}

	wg.Wait()
}
