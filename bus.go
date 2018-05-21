package eventsource

type EventBus interface {
	EventPublisher
	EventSubscriber
}

type EventPublisher interface {
	Publish(exchange string, topic string, event Event) error
}

type EventSubscriber interface {
	Subscribe(topic string, event Event)
}
