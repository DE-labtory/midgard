package eventsource_test

import (
	"fmt"
	"log"
	"testing"

	"reflect"

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

	typ := reflect.TypeOf(cmd)
	fmt.Println(typ.Name())

	b, err := json.Marshal(cmd)
	assert.NoError(t, err)

	fmt.Println(b)

	//typ := reflect.TypeOf(cmd)
	//v := reflect.New(typ)
	//
	//initializeStruct(typ, v.Elem())
	//paramInterface := v.Interface()
	//
	////cmd2 := UserNameUpdateCommand{}
	////
	////fmt.Println(reflect.TypeOf(cmd2))
	////fmt.Println(reflect.TypeOf(paramInterface))
	////
	//err = json.Unmarshal(b, paramInterface)
	//fmt.Println(reflect.ValueOf(paramInterface).Elem())
	////err = json.Unmarshal(b, &cmd2)
	////assert.NoError(t, err)
	//fmt.Println(reflect.ValueOf(paramInterface).Elem().Interface())
	////err = d.Route(b, "UserNameUpdateCommand")

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
	fmt.Println("hello world2")
}

//func TestMarshal(t *testing.T) {
//
//	cmd := UserNameUpdateCommand{
//		Name: "jun",
//	}
//
//	typ := reflect.TypeOf(cmd)
//	v := reflect.New(typ)
//	initializeStruct(typ, v.Elem())
//	emtpy := v.Interface()
//
//	b, _ := json.Marshal(cmd)
//
//	json.Unmarshal(b, emtpy)
//
//	fmt.Println(emtpy)
//
//	command := emtpy.(eventsource.Command)
//
//	//method.Func.Call([]reflect.Value{sourceValue, eventValue})
//}

func call(user UserNameUpdateCommand) {

}

func initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}
