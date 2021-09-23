package rabbitmq

import (
	"github.com/streadway/amqp"
)

func (c *Client) initQueues(queues []string) error {
	for _, queue := range queues {
		err := c.CreateQueue(queue)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewClient(serviceName, URL string, queues []string) (*Client, error) {
	config, err := newConfig()
	if err != nil {
		return nil, err
	}
	connect, err := amqp.Dial(URL)
	if err != nil {
		return nil, err
	}
	channel, err := connect.Channel()
	if err != nil {
		return nil, err
	}
	var client = &Client{
		name:       serviceName,
		connection: connect,
		channel:    channel,
		config:     config,
	}
	err = client.initQueues(queues)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) SendMessage(from, to, msg, ID string, body []byte) error {
	queue := to + "." + msg
	replyTo := from + "." + msg
	return c.channel.Publish("", queue, c.config.Mandatory, c.config.Immediate, amqp.Publishing{
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
		c.config.Channel.AutoAck,
		c.config.Channel.Exclusive,
		c.config.Channel.NoLocal,
		c.config.Channel.NoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}
	m := <-message
	return NewMessage(m.Type, m.CorrelationId, m.ReplyTo, m.Body), nil
}

func (c *Client) CreateQueue(name string) error {
	queueName := c.name + "." + name
	_, err := c.channel.QueueDeclare(
		queueName,
		c.config.Queue.Durable,
		c.config.Queue.AutoDel,
		c.config.Queue.Exclusive,
		c.config.Queue.NoWait,
		nil,
	)
	return err
}
