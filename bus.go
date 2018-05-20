package eventsource

type EventBus interface {
	Publisher
	Subscriber
}

type Publisher interface {
	Publish(topic string, event Event)
}

type Subscriber interface {
	Subscribe(topic string)
}
