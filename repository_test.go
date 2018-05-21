package eventsource_test

import (
	"fmt"
	"testing"

	"errors"

	"os"

	"github.com/it-chain/eventsource"
	"github.com/it-chain/eventsource/store/leveldb"
	"github.com/stretchr/testify/assert"
)

//aggregate
type UserAggregate struct {
	Name string
	eventsource.EventModel
}

func (u *UserAggregate) On(event eventsource.Event) error {

	switch v := event.(type) {

	case *UserCreated:
		u.AggregateID = v.AggregateID

	case *UserUpdated:
		u.Name = v.Name

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

//event
type UserCreated struct {
	eventsource.EventModel
}

//event
type UserUpdated struct {
	Name string
	eventsource.EventModel
}

func TestNewRepository(t *testing.T) {

	path := "test"
	defer os.RemoveAll(path)

	store := leveldb.NewEventStore(path, leveldb.NewSerializer(UserCreated{}, UserUpdated{}))
	r := eventsource.NewRepo(store, nil)

	aggregateID := "123"

	err := r.Save(aggregateID, UserCreated{
		EventModel: eventsource.EventModel{
			AggregateID: aggregateID,
			Type:        "User",
		},
	})

	assert.NoError(t, err)

	err = r.Save(aggregateID, UserUpdated{
		EventModel: eventsource.EventModel{
			AggregateID: aggregateID,
			Type:        "User",
		},
		Name: "jun",
	})

	assert.NoError(t, err)

	err = r.Save(aggregateID, UserUpdated{
		EventModel: eventsource.EventModel{
			AggregateID: aggregateID,
			Type:        "User",
		},
		Name: "jun2",
	})

	assert.NoError(t, err)

	user := &UserAggregate{}

	err = r.Load(user, aggregateID)
	assert.NoError(t, err)

	assert.Equal(t, user.AggregateID, aggregateID)
	assert.Equal(t, user.Name, "jun2")

	fmt.Println(user)
}
