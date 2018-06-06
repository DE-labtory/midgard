package mongodb

import (
	"errors"
	"testing"
	"time"

	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/store"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserAddedEvent struct {
	Name string
	midgard.EventModel
}

func TestNewEventStore(t *testing.T) {
	// given
	path := "mongodb://localhost:27017"
	dbname := "test"

	defer dropDB(path, dbname)

	// When
	store, err := NewEventStore(path, dbname, store.NewSerializer(UserAddedEvent{}))

	// then
	assert.NotNil(t, store)
	assert.Nil(t, err)

}

func TestNewEventStore_WrongPath(t *testing.T) {
	// given
	wrongpath := "strange_path"
	dbname := "test"

	// When
	_, err := NewEventStore(wrongpath, dbname, store.NewSerializer(UserAddedEvent{}))

	// Then
	assert.Error(t, err)
}

func TestStore_Save(t *testing.T) {
	// given
	path := "mongodb://localhost:27017"
	dbname := "test"
	session, _ := mgo.Dial(path)

	defer func() {
		dropDB(path, dbname)
		session.Close()
	}()

	store, err := NewEventStore(path, dbname, store.NewSerializer(UserAddedEvent{}))
	assert.NoError(t, err)

	document := Document{}
	var aggregateID string
	aggregateID = "1"

	events := []UserAddedEvent{
		{Name: "zf1", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 1}},
		{Name: "zf2", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 1}},
		{Name: "zf3", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 1}},
	}

	// When
	saveErr := store.Save(aggregateID, ToEvent(events...)...)

	// Then
	assert.Equal(t, nil, saveErr)

	// When
	session.DB(dbname).C("events").Find(bson.M{"aggregate_id": aggregateID}).One(&document)

	// Then
	assert.Equal(t, 3, len(document.History))
	assert.Equal(t, "1", document.AggregateID)

	// When
	events2 := []UserAddedEvent{
		{Name: "jun", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 2}},
		{Name: "jun2", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 2}},
	}
	saveErr2 := store.Save(aggregateID, ToEvent(events2...)...)

	// Then
	assert.Equal(t, nil, saveErr2)

	// When
	session.DB(dbname).C("events").Find(bson.M{"aggregate_id": aggregateID}).One(&document)

	// Then
	assert.Equal(t, 5, len(document.History))
	assert.Equal(t, "1", document.AggregateID)
}

func TestStore_Load(t *testing.T) {
	// given
	path := "mongodb://localhost:27017"
	dbname := "test"

	defer dropDB(path, dbname)

	store, err := NewEventStore(path, dbname, store.NewSerializer(UserAddedEvent{}))
	assert.NoError(t, err)

	var aggregateID string
	aggregateID = "1"

	events := []UserAddedEvent{
		{Name: "zf1", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 1}},
		{Name: "zf2", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 1}},
		{Name: "zf3", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 1}},
	}

	// When
	store.Save(aggregateID, ToEvent(events...)...)
	Events, err := store.Load(aggregateID)

	// Then
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(Events))
	assert.Equal(t, aggregateID, Events[0].GetID())

	userAddedEvent := Events[0].(*UserAddedEvent)
	assert.Equal(t, "zf1", userAddedEvent.Name)
}

func TestStore_Load_NoMatchingDocument(t *testing.T) {
	// given
	path := "mongodb://localhost:27017"
	dbname := "test"

	defer dropDB(path, dbname)

	store, err := NewEventStore(path, dbname, store.NewSerializer(UserAddedEvent{}))
	assert.NoError(t, err)

	var aggregateID string
	aggregateID = "1"

	// When
	_, err = store.Load(aggregateID)
	// Then
	assert.Equal(t, errors.New("not found"), err)
}

func dropDB(path string, dbname string) {
	session, _ := mgo.Dial(path)

	defer session.Close()

	err := session.DB(dbname).DropDatabase()

	if err != nil {
		panic(err)
	}
}

// Convert a slice or array of a specific type to array of midgard.Event
func ToEvent(event ...UserAddedEvent) []midgard.Event {
	intf := make([]midgard.Event, len(event))
	for i, v := range event {

		intf[i] = midgard.Event(v)
	}
	return intf
}
