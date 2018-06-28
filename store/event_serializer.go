package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/it-chain/midgard"
)

type SerializedEvent struct {
	Type string
	Data []byte
}

type EventSerializer interface {
	// MarshalEvent converts an Event to a Record
	Marshal(event midgard.Event) (SerializedEvent, error)

	// UnmarshalEvent converts an Event backed into a Record
	Unmarshal(serializedEvent SerializedEvent) (midgard.Event, error)

	// RegisterEvent
	Register(events ...midgard.Event)
}

type JSONSerializer struct {
	eventTypes map[string]reflect.Type
}

func NewSerializer(events ...midgard.Event) EventSerializer {

	s := &JSONSerializer{
		eventTypes: make(map[string]reflect.Type),
	}

	s.Register(events...)

	return s
}

func (j *JSONSerializer) Register(events ...midgard.Event) {

	for _, event := range events {
		rawType, name := GetTypeName(event)
		j.eventTypes[name] = rawType
	}
}

func (j *JSONSerializer) Marshal(e midgard.Event) (SerializedEvent, error) {

	serializedEvent := SerializedEvent{}
	_, name := GetTypeName(e)
	serializedEvent.Type = name

	data, err := json.Marshal(e)

	if err != nil {
		return SerializedEvent{}, err
	}

	serializedEvent.Data = data

	return serializedEvent, nil
}

func (j *JSONSerializer) Unmarshal(serializedEvent SerializedEvent) (midgard.Event, error) {

	t, ok := j.eventTypes[serializedEvent.Type]

	if !ok {
		return nil, errors.New(fmt.Sprintf("unbound event type, %v", serializedEvent.Type))
	}

	v := reflect.New(t).Interface()

	err := json.Unmarshal(serializedEvent.Data, v)
	if err != nil {
		return nil, err
	}

	return v.(midgard.Event), nil
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
