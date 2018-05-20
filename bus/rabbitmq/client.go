package rabbitmq

import (
	"encoding/json"

	"github.com/it-chain/eventsource"
	"github.com/streadway/amqp"
)

type Client struct {
	conn *amqp.Connection
}

func Connect(rabbitmqUrl string) *Client {

	if rabbitmqUrl == "" {
		rabbitmqUrl = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitmqUrl)

	if err != nil {
		panic("Failed to connect to RabbitMQ" + err.Error())
	}

	return &Client{
		conn: conn,
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Publish(exchange string, topic string, event eventsource.Event) error {

	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return err
	}

	data, err := json.Marshal(event)

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
			Body:        []byte(data),
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

	defer ch.Close()

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

type Handler func(delivery amqp.Delivery)

func (c *Client) Consume(exchange string, topic string, handler Handler) error {

	chanDelivery, err := c.consume(exchange, topic)

	if err != nil {
		return err
	}

	go func() {
		for delivery := range chanDelivery {
			handler(delivery)
		}
	}()

	return nil
}
