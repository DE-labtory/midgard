package eventsource

import (
	"errors"
	"reflect"
	"strings"

	"fmt"

	"github.com/altairsix/eventsource"
)

type Dispatcher interface {
	Dispatch(command Command)
}

type HandlerFunc interface{}

type CommandDispatcher struct {
	handlerMap map[reflect.Type]func(command Command)
}

func NewDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		handlerMap: make(map[reflect.Type]func(command Command)),
	}
}

func (c *CommandDispatcher) SetHandler(handler interface{}) error {

	if reflect.TypeOf(handler).Kind() != reflect.Ptr {
		return errors.New("handler should be ptr type")
	}

	sourceType := reflect.TypeOf(handler)
	methodCount := sourceType.NumMethod()

	for i := 0; i < methodCount; i++ {
		method := sourceType.Method(i)

		if method.Type.NumIn() != 2 {
			return errors.New("number of parameter of handler is not 2")
		}

		commandType := method.Type.In(1)
		commandInterface := reflect.TypeOf((*eventsource.Command)(nil)).Elem()

		if !commandType.ConvertibleTo(commandInterface) {
			return errors.New(fmt.Sprintf("method parameter does not implementing command [%s]", commandType))
		}

		handler := createEventHandler(method, handler)
		c.handlerMap[commandType] = handler
	}

	return nil
}

func createEventHandler(method reflect.Method, handler interface{}) func(Command) {
	return func(command Command) {
		sourceValue := reflect.ValueOf(handler)
		eventValue := reflect.ValueOf(command)

		// Call actual event handling method.
		method.Func.Call([]reflect.Value{sourceValue, eventValue})
	}
}

func (c CommandDispatcher) Dispatch(command eventsource.Command) error {
	eventType := reflect.TypeOf(command)
	if handler, ok := c.handlerMap[eventType]; ok {
		handler(command)
	} else {
		return errors.New(fmt.Sprintf("No handler found for event: %v in %v", eventType.String()))
	}

	return nil
}

func IsFunc(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Func
}

func GetTypeName(source interface{}) (reflect.Type, string) {

	rawType := reflect.TypeOf(source)

	if rawType.Kind() == reflect.Ptr {
		rawType = rawType.Elem()
	}

	name := rawType.String()
	parts := strings.Split(name, ".")
	return rawType, parts[1]
}
