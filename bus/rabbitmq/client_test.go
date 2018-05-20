package rabbitmq_test

import (
	"fmt"
	"log"
	"testing"

	"sync"

	"time"

	"github.com/it-chain/eventsource"
	"github.com/it-chain/eventsource/bus/rabbitmq"
	"github.com/stretchr/testify/assert"
)

var wg sync.WaitGroup

func TestConnect(t *testing.T) {

	wg.Add(1)
	c := rabbitmq.Connect("")
	err := c.Consume("asd", "asd", &Dispatcher{})
	assert.NoError(t, err)

	err = c.Publish("asd", "asd", UserNameUpdateCommand{
		Name: "JUN",
		EventModel: eventsource.EventModel{
			AggregateID: "123",
			Time:        time.Now(),
			Type:        "123",
			Version:     1,
		}})
	assert.NoError(t, err)

	wg.Wait()
}

type UserNameUpdateCommand struct {
	eventsource.EventModel
	Name string
}

type UserAddCommand struct {
	eventsource.EventModel
}

type Dispatcher struct {
}

func (d *Dispatcher) Handle(command UserAddCommand) {
	log.Print("hello world")
}

func (d *Dispatcher) HandleNameUpdateCommand(command UserNameUpdateCommand) {
	fmt.Println("hello world2")
	fmt.Println(command)
	wg.Done()
}
