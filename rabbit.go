package rabbitmq

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
		name:       serviceName,
		connection: connect,
		channel:    channel,
	}, nil
}

func (c *Client) SendMessage(from, to, msg, ID string, body []byte) error {
	queue := to + "." + msg
	replyTo := from + "." + msg
	return c.channel.Publish("", queue, false, false, amqp.Publishing{
		ContentType:   "application/json",
		Type:          msg,
		CorrelationId: ID,
		ReplyTo:       replyTo,
		Body:          body,
	})
}

func (c *Client) ReceiveMessage(msg string) (*Message, error) {
	queue := c.name + "." + msg
	message, err := c.channel.Consume(
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
	m := <-message
	return NewMessage(m.Type, m.CorrelationId, m.ReplyTo, m.Body), nil
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
