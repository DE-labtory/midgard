package main

import (
	"os"

	"github.com/it-chain/eventsource"
	"github.com/it-chain/eventsource/store/leveldb"
)

type User struct {
	name string
	eventsource.Event
}

func main() {
	path := "test"
	store := leveldb.NewEventStore(path)
	defer os.RemoveAll(path)

	store.Save("1", User{})

	_, err := store.Load("1")

	if err != nil {
		panic(err.Error())
	}

}
