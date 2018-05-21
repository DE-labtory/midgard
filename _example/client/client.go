package main

import (
	"github.com/it-chain/eventsource"
	"github.com/it-chain/eventsource/bus/rabbitmq"
)

type UserCreateCommand struct {
	eventsource.CommandModel
}

type UserNameUpdateCommand struct {
	eventsource.CommandModel
	Name string
}

func main() {

	c := rabbitmq.Connect("")

	err := c.Publish("Command", "User", UserCreateCommand{})

	if err != nil {
		panic(err.Error())
	}

	err = c.Publish("Command", "User", UserNameUpdateCommand{

		Name: "jun2",
		CommandModel: eventsource.CommandModel{
			AggregateID: "123",
		},
	})

	if err != nil {
		panic(err.Error())
	}
}
