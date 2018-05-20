package eventsource

import (
	"errors"
	"fmt"
	"reflect"
)

type Router interface {
	Route(command Command)
}

type ParamBasedRouter struct {
	handlerMap map[reflect.Type]func(param interface{})
}

func NewParamBasedRouter(handlers ...interface{}) (*ParamBasedRouter, error) {

	p := &ParamBasedRouter{
		handlerMap: make(map[reflect.Type]func(param interface{})),
	}

	for _, handler := range handlers {
		err := p.SetHandler(handler)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (c *ParamBasedRouter) SetHandler(handler interface{}) error {

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

		paramType := method.Type.In(1)

		_, ok := c.handlerMap[paramType]

		if ok {
			return errors.New(fmt.Sprintf("same param already exist [%s]", paramType))
		}

		handler := createEventHandler(method, handler)
		c.handlerMap[paramType] = handler
	}

	return nil
}

func createEventHandler(method reflect.Method, handler interface{}) func(interface{}) {
	return func(param interface{}) {
		sourceValue := reflect.ValueOf(handler)
		eventValue := reflect.ValueOf(param)

		// Call actual event handling method.
		method.Func.Call([]reflect.Value{sourceValue, eventValue})
	}
}

func (c ParamBasedRouter) Dispatch(param interface{}) error {

	eventType := reflect.TypeOf(param)
	if handler, ok := c.handlerMap[eventType]; ok {
		handler(param)
	} else {
		return errors.New(fmt.Sprintf("No handler found for command [%s]", eventType.String()))
	}

	return nil
}
