package bus

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var ErrTypeNotFound = errors.New("Type of handler not found")

//Route data depends on type
type Router interface {

	//route data depends on matching value
	Route(data []byte, matchingValue string) error

	SetHandler(handler interface{}) error
}

//ParamBasedRouter routes data through the structure and structure name(matching value) of the parameter
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

////handler should be a struct pointer which has handler method
func (c *ParamBasedRouter) SetHandler(handler interface{}) error {

	if reflect.TypeOf(handler).Kind() != reflect.Ptr {
		return errors.New("handler should be ptr type")
	}

	sourceType := reflect.TypeOf(handler)
	methodCount := sourceType.NumMethod() // NumMethod returns the number of exported methods in the value's method set.

	// executes every method if handler interface
	for i := 0; i < methodCount; i++ {

		// Method returns a function value corresponding to v's i'th method.
		// The arguments to a Call on the returned function should not include a receiver;
		// the returned function will always use v as the receiver.
		// Method panics if i is out of range or if v is a nil interface value.
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

func (c ParamBasedRouter) Route(data []byte, structName string) (err error) {

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	paramType, handler, err := c.findTypeOfHandler(structName)

	if err != nil {
		errors.New(fmt.Sprintf("No handler found for struct [%s]", structName))
	}

	v := reflect.New(paramType)
	initializeStruct(paramType, v.Elem())
	paramInterface := v.Interface()

	err = json.Unmarshal(data, paramInterface)

	if err != nil {
		return err
	}

	paramValue := reflect.ValueOf(paramInterface).Elem().Interface()

	handler(paramValue) //execute handler

	return nil
}

//find type of handler by struct name
func (c ParamBasedRouter) findTypeOfHandler(typeName string) (reflect.Type, func(param interface{}), error) {

	for paramType, handler := range c.handlerMap {
		name := paramType.Name()

		if name == typeName {
			return paramType, handler, nil
		}
	}

	return nil, nil, ErrTypeNotFound
}

//build empty struct from struct type
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
