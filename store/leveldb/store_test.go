package leveldb_test

import (
	"testing"

	"time"

	"os"

	"github.com/it-chain/eventsource"
	"github.com/it-chain/eventsource/store/leveldb"
	"github.com/stretchr/testify/assert"
)

type UserAddedEvent struct {
	Name string
	eventsource.EventModel
}

func TestNew(t *testing.T) {

	//given
	path := "test"
	store := leveldb.NewEventStore(path, leveldb.NewSerializer(UserAddedEvent{}))
	defer os.RemoveAll(path)

	var aggregateID string

	aggregateID = "1"

	events := []UserAddedEvent{
		{Name: "jun", EventModel: eventsource.EventModel{AggregateID: aggregateID, Time: time.Now().UTC(), Version: 1}},
		{Name: "jun2", EventModel: eventsource.EventModel{AggregateID: aggregateID, Time: time.Now().UTC(), Version: 1}},
		{Name: "jun3", EventModel: eventsource.EventModel{AggregateID: aggregateID, Time: time.Now().UTC(), Version: 1}},
	}

	events2 := []UserAddedEvent{
		{Name: "jun", EventModel: eventsource.EventModel{AggregateID: aggregateID, Time: time.Now().UTC(), Version: 2}},
		{Name: "jun2", EventModel: eventsource.EventModel{AggregateID: aggregateID, Time: time.Now().UTC(), Version: 2}},
	}

	//when
	err := store.Save(aggregateID, ToEvent(events...)...)
	assert.NoError(t, err)

	history, err := store.Load(aggregateID)
	assert.NoError(t, err)

	//then
	assert.Equal(t, ToUserAddedEvent(t, history...), events)

	//when
	store.Save(aggregateID, ToEvent(events2...)...)
	assert.NoError(t, err)

	totalEvents := append(events, events2...)

	history, err = store.Load(aggregateID)
	assert.NoError(t, err)

	//then
	assert.Equal(t, ToUserAddedEvent(t, history...), totalEvents)
}

////Convert a slice or array of a specific type to array of interface{}
func ToEvent(event ...UserAddedEvent) []eventsource.Event {
	intf := make([]eventsource.Event, len(event))
	for i, v := range event {
		intf[i] = eventsource.Event(v)
	}
	return intf
}

func ToUserAddedEvent(t *testing.T, events ...eventsource.Event) []UserAddedEvent {

	uae := make([]UserAddedEvent, 0)
	for _, v := range events {
		userAddedEvent := v.(*UserAddedEvent)
		uae = append(uae, *userAddedEvent)
	}

	return uae
}
