package main

import (
	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
)

type UserCreateCommand struct {
	midgard.CommandModel
}

type UserNameUpdateCommand struct {
	midgard.CommandModel
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
		CommandModel: midgard.CommandModel{
			AggregateID: "123",
		},
	})

	if err != nil {
		panic(err.Error())
	}
}
