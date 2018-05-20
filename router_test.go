package eventsource_test

import (
	"log"
	"testing"

	"encoding/json"

	"github.com/it-chain/eventsource"
	"github.com/stretchr/testify/assert"
)

func TestNewParamBasedRouter(t *testing.T) {
	d, err := eventsource.NewParamBasedRouter()
	assert.NoError(t, err)

	err = d.SetHandler(&Dispatcher{})
	assert.NoError(t, err)

	cmd := UserNameUpdateCommand{
		Name: "jun",
	}

	b, _ := json.Marshal(cmd)

	err = d.Route(b, "UserNameUpdateCommand")
	assert.NoError(t, err)
}

type UserNameUpdateCommand struct {
	eventsource.CommandModel
	Name string
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
