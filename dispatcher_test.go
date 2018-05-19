package eventsource_test

import (
	"testing"

	"log"

	"github.com/it-chain/eventsource"
	"github.com/stretchr/testify/assert"
)

func TestDefaultDispatcher_SetHandler(t *testing.T) {
	d := eventsource.NewDispatcher()
	err := d.SetHandler(&Dispatcher{})
	assert.NoError(t, err)

	d.Dispatch(UserAddCommand{})
	d.Dispatch(UserNameUpdateCommand{})
	//sourceType, name := eventsource.GetTypeName(&Dispatcher{})
	//log.Print(sourceType.NumMethod())
	//log.Print(name)
	//
	//fooType := reflect.TypeOf(&Dispatcher{})
	//log.Print(fooType.NumMethod())
	//
	//for i := 0; i < fooType.NumMethod(); i++ {
	//	method := fooType.Method(i)
	//	fmt.Println(method.Name)
	//}
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

//type Dispatcher struct {
//}
//
//func (f *Dispatcher) Hanlde(asd string) {
//
//}
//
//func (f *Dispatcher) Baz() {
//}
