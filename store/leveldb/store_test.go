package leveldb_test

import (
	"testing"

	"time"

	"reflect"

	"fmt"

	"os"

	"github.com/it-chain/eventsource"
	"github.com/it-chain/eventsource/store/leveldb"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	//given
	path := "test"
	store := leveldb.NewEventStore(path)
	defer os.RemoveAll(path)

	var aggregateID string

	aggregateID = "1"

	events := []eventsource.EventD{{AggregateID: aggregateID, Time: time.Now().UTC(), Version: 1}, {AggregateID: aggregateID, Time: time.Now().UTC(), Version: 1}, {AggregateID: aggregateID, Time: time.Now().UTC(), Version: 1}}
	events2 := []eventsource.EventD{{AggregateID: aggregateID, Time: time.Now().UTC(), Version: 2}, {AggregateID: aggregateID, Time: time.Now().UTC(), Version: 2}}

	//when
	err := store.Save(aggregateID, ToEvent(events...)...)
	assert.NoError(t, err)

	history, err := store.Load(aggregateID)
	assert.NoError(t, err)

	//then
	assert.Equal(t, history, events)

	//when
	store.Save(aggregateID, ToEvent(events2...)...)
	assert.NoError(t, err)

	totalEvents := append(events, events2...)

	history, err = store.Load(aggregateID)
	assert.NoError(t, err)

	//then
	assert.Equal(t, totalEvents, history)
}

func TestNewEventStore(t *testing.T) {
	var aggregateID string

	aggregateID = "1"
	events2 := []eventsource.EventD{{AggregateID: aggregateID, Time: time.Now().UTC(), Version: 2}, {AggregateID: aggregateID, Time: time.Now().UTC(), Version: 2}}
	v := reflect.ValueOf(events2)

	fmt.Print(v)
}

//
//Convert a slice or array of a specific type to array of interface{}
func ToEvent(eventD ...eventsource.EventD) []eventsource.Event {
	intf := make([]eventsource.Event, len(eventD))
	for i, v := range eventD {
		intf[i] = eventsource.Event(v)
	}
	return intf
}
