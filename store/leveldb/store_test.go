package leveldb_test

import (
	"testing"

	"os"

	"time"

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

	events := []eventsource.Event{{AggregateID: aggregateID, Time: time.Now().UTC(), Version: 1}, {AggregateID: aggregateID, Time: time.Now().UTC(), Version: 1}, {AggregateID: aggregateID, Time: time.Now().UTC(), Version: 1}}
	events2 := []eventsource.Event{{AggregateID: aggregateID, Time: time.Now().UTC(), Version: 2}, {AggregateID: aggregateID, Time: time.Now().UTC(), Version: 2}}

	//when
	err := store.Save(aggregateID, events...)
	assert.NoError(t, err)

	history, err := store.Load(aggregateID)
	assert.NoError(t, err)

	//then
	assert.Equal(t, history, events)

	//when
	store.Save(aggregateID, events2...)
	assert.NoError(t, err)

	totalEvents := append(events, events2...)

	history, err = store.Load(aggregateID)
	assert.NoError(t, err)

	//then
	assert.Equal(t, totalEvents, history)
}
