package mongodb

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/it-chain/midgard"
	"gopkg.in/mgo.v2"
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
	store := NewEventStore(path, dbname)

	// then
	assert.NotEqual(t, store, nil)

}

func TestNewEventStore_WrongPath(t *testing.T) {
	// given
	wrongpath := "strange_path"
	dbname := "test"

	defer dropDB(wrongpath, dbname)

	// When
	store := NewEventStore(wrongpath, dbname)

	// Then
	assert.Equal(t, store, nil)
}


//func TestStore_Save(t *testing.T) {
//	// given
//	path := "mongodb://localhost:27017"
//	dbname := "test"
//
//	defer dropDB(path, dbname)
//
//	session, _ := mgo.Dial(path)
//	store := NewEventStore(path, dbname)
//
//	user := []midgard.Event{}
//	var aggregateID string
//	aggregateID = "1"
//
//	events := []UserAddedEvent{
//		{Name: "zf1", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 1}},
//		{Name: "zf2", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 1}},
//		{Name: "zf3", EventModel: midgard.EventModel{ID: aggregateID, Time: time.Now().UTC(), Version: 1}},
//	}
//
//	// When
//	saveErr := store.Save(aggregateID, ToEvent(events...)...)
//
//	// Then
//	assert.Equal(t, saveErr, nil)
//
//	// When
//	session.DB(dbname).C("events").Find(bson.M{"Name": "zf1"}).One(&user)
//	fmt.Println(user)
//	//u := midgard.EventModel(users[0])
//	//fmt.Println(user.GetID())
//}


func dropDB(path string, dbname string) {
	session, _ := mgo.Dial(path)

	defer session.Close()

	err := session.DB(dbname).DropDatabase()

	if err != nil {
		panic(err)
	}
}
//
//// Convert a slice or array of a specific type to array of midgard.Event
//func ToEvent(event ...UserAddedEvent) []midgard.Event {
//	intf := make([]midgard.Event, len(event))
//	for i, v := range event {
//
//		intf[i] = midgard.Event(v)
//	}
//	return intf
//}
//
//func ToUserAddedEvent(events ...midgard.Event) []UserAddedEvent {
//	uae := make([]UserAddedEvent, 0)
//	for _, v := range events {
//		userAddedEvent := v.(*UserAddedEvent)
//		uae = append(uae, *userAddedEvent)
//	}
//
//	return uae
//}

