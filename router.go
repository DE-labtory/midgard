package eventsource

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var ErrTypeNotFound = errors.New("Type of handler not found")

type Router interface {
	Route(data []byte, typeName string) error
	SetHandler(handler interface{}) error
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

func (c ParamBasedRouter) Route(data []byte, typeName string) error {

	paramType, handler, err := c.findTypeOfHandler(typeName)

	if err != nil {
		errors.New(fmt.Sprintf("No handler found for param [%s]", typeName))
	}

	v := reflect.New(paramType)
	initializeStruct(paramType, v.Elem())
	paramInterface := v.Interface()

	err = json.Unmarshal(data, paramInterface)

	if err != nil {
		return err
	}

	paramValue := reflect.ValueOf(paramInterface).Elem().Interface()

	handler(paramValue)

	return nil
}

func (c ParamBasedRouter) findTypeOfHandler(typeName string) (reflect.Type, func(param interface{}), error) {

	for paramType, handler := range c.handlerMap {
		name := paramType.Name()

		if name == typeName {
			return paramType, handler, nil
		}
	}

	return nil, nil, ErrTypeNotFound
}

func initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)

		if !f.CanSet() {
			continue
		}

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
