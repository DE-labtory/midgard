package event_store_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/event_store"
	"github.com/stretchr/testify/assert"
	mgo "gopkg.in/mgo.v2"
)

type MockPublisher struct {
	publishFunc func(exchange string, topic string, data interface{}) (err error)
}

func (m MockPublisher) Publish(exchange string, topic string, data interface{}) (err error) {
	return m.publishFunc(exchange, topic, data)
}

type UserCreatedEvent struct {
	midgard.EventModel
}

type UserNameUpdatedEvent struct {
	midgard.EventModel
	Name string
}

func TestSave(t *testing.T) {

	path := "mongodb://localhost:27017"
	dbname := "test"

	m := MockPublisher{}
	m.publishFunc = func(exchange string, topic string, data interface{}) (err error) {

		return nil
	}

	event_store.InitMongoEventStore(path, dbname, m)

	defer dropDB(path, dbname)

	err := event_store.Save("123", UserCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   "123",
			Type: "User",
		},
	})

	assert.NoError(t, err)
}

// aggregate
type User struct {
	Name string
	midgard.AggregateModel
}

func (u *User) On(event midgard.Event) error {

	switch v := event.(type) {

	case *UserCreatedEvent:
		u.ID = v.ID

	case *UserNameUpdatedEvent:
		u.Name = v.Name

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

func TestLoad(t *testing.T) {

	path := "mongodb://localhost:27017"
	dbname := "test"

	m := MockPublisher{}
	m.publishFunc = func(exchange string, topic string, data interface{}) (err error) {

		return nil
	}

	event_store.InitMongoEventStore(path, dbname, m)
	event_store.RegisterEvents(UserCreatedEvent{})

	defer dropDB(path, dbname)

	err := event_store.Save("123", UserCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   "123",
			Type: "User",
		},
	})

	assert.NoError(t, err)

	user := &User{}
	err = event_store.Load(user, "123")
	assert.NoError(t, err)

	assert.Equal(t, "123", user.AggregateModel.ID)
}

func dropDB(path string, dbname string) {
	session, _ := mgo.Dial(path)

	defer session.Close()

	err := session.DB(dbname).DropDatabase()

	if err != nil {
		panic(err)
	}
}
