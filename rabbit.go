package RabbitMQ

import (
	"github.com/streadway/amqp"
)

func NewClient(serviceName, URL string) (*Client, error) {
	connect, err := amqp.Dial(URL)
	if err != nil {
		return nil, err
	}
	channel, err := connect.Channel()
	if err != nil {
		return nil, err
	}
	return &Client{
		name:    serviceName,
		channel: channel,
	}, nil
}

func (c *Client) SendMessage(from, to, msg, ID, body string) error {
	queue := to + "." + msg
	replyTo := from + "." + msg
	return c.channel.Publish("", queue, false, false, amqp.Publishing{
		ContentType:   "application/json",
		Type:          msg,
		CorrelationId: ID,
		ReplyTo:       replyTo,
		Body:          []byte(body),
	})
}

func (c *Client) ReceiveMessage(msg string) (*Message, error) {
	queue := c.name + "." + msg
	messages, err := c.channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return NewMessage((<-messages).Type, (<-messages).CorrelationId, (<-messages).ReplyTo, (<-messages).Body), nil
}

func (c *Client) CreateQueue(message string, isResponse bool) error {
	queueName := c.name + "." + message
	if isResponse {
		queueName += ".get_response"
	}
	_, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		true,
		false,
		nil,
	)
	return err
}
