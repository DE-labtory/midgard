package midgard

type Bus interface {
	Publisher
	Subscriber
}

type Publisher interface {
	Publish(topic string, data interface{}) (err error)
}

type Subscriber interface {
	Subscribe(topic string, source interface{}) error
}
