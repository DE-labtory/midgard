package eventsource_test

import (
	"log"
	"testing"

	"github.com/it-chain/eventsource"
	"github.com/stretchr/testify/assert"
)

func TestNewParamBasedRouter(t *testing.T) {
	d, err := eventsource.NewParamBasedRouter()
	assert.NoError(t, err)

	err = d.SetHandler(&Dispatcher{})
	assert.NoError(t, err)

	d.Dispatch(UserAddCommand{})
	d.Dispatch(UserNameUpdateCommand{})
	assert.NoError(t, err)
}

type UserNameUpdateCommand struct {
	eventsource.CommandModel
	name string
}

type UserAddCommand struct {
	eventsource.CommandModel
}

type Dispatcher struct {
}

func (d *Dispatcher) Handle(command UserAddCommand) {
	log.Print("hello world")
}

func (d *Dispatcher) HandleNameUpdateCommand(command UserNameUpdateCommand) {
	log.Print("hello world2")
}
