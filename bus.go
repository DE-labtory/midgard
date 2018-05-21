package eventsource

type Bus interface {
	Publisher
	Subscriber
}

type Publisher interface {
	Publish(exchange string, topic string, data interface{}) (err error)
}

type Subscriber interface {
	Subscribe(exchange string, topic string, source interface{}) error
}
