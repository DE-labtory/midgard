package rabbitmq

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"

	"github.com/it-chain/midgard/bus"
	"github.com/streadway/amqp"
)

type Message struct {
	MatchingValue string
	Data          []byte
}

type Client struct {
	conn   *amqp.Connection
	router bus.Router
}

func Connect(rabbitmqUrl string) *Client {

	if rabbitmqUrl == "" {
		rabbitmqUrl = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitmqUrl)

	if err != nil {
		panic("Failed to connect to RabbitMQ" + err.Error())
	}

	p, _ := bus.NewParamBasedRouter()

	return &Client{
		conn:   conn,
		router: p,
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Publish(exchange string, topic string, data interface{}) (err error) {

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	ch, err := c.conn.Channel()

	if err != nil {
		return err
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return err
	}

	var matchingValue string

	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		matchingValue = reflect.TypeOf(data).Elem().Name()
	} else {
		matchingValue = reflect.TypeOf(data).Name()
	}

	b, err := json.Marshal(data)

	message := Message{
		MatchingValue: matchingValue,
		Data:          b,
	}

	byte, err := json.Marshal(message)

	if err != nil {
		return err
	}

	if err != nil {
		panic("Failed to open exchange" + err.Error())
	}

	err = ch.Publish(
		exchange, // exchange
		topic,    // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        byte,
		})

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) consume(exchange string, topic string) (<-chan amqp.Delivery, error) {

	ch, err := c.conn.Channel()

	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		panic("Failed to open a channel" + err.Error())
	}

	err = ch.QueueBind(
		q.Name,   // queue name
		topic,    // routing key
		exchange, // exchange
		false,
		nil)

	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (c *Client) Subscribe(exchange string, topic string, source interface{}) error {

	chanDelivery, err := c.consume(exchange, topic)

	if err != nil {
		return err
	}

	err = c.router.SetHandler(source)

	if err != nil {
		return err
	}

	go func() {
		for delivery := range chanDelivery {
			data := delivery.Body

			message := &Message{}
			err := json.Unmarshal(data, message)

			if err != nil {
				log.Print(err.Error())
			}

			c.router.Route(message.Data, message.MatchingValue)
		}
	}()

	return nil
}
