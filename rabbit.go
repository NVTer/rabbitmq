package rabbitmq

import (
	"github.com/streadway/amqp"
)

func NewClient(serviceName, URL string) (*Client, error) {
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
	return &Client{
		name:       serviceName,
		connection: connect,
		channel:    channel,
		config:     config,
	}, nil
}

func (c *Client) SendMessage(msgType, ID string, body []byte, isResponse bool) error {
	var replyTo string
	if isResponse {
		replyTo = c.name + "." + msgType + "_response"
	}
	return c.channel.Publish(msgType, "", c.config.Mandatory, c.config.Immediate, amqp.Publishing{
		ContentType:   "application/json",
		Type:          msgType,
		CorrelationId: ID,
		ReplyTo:       replyTo,
		Body:          body,
	})
}

func (c *Client) StartConsumer(msgType string, msgChannel chan Message, isResponse bool) error {
	queue := c.name + "." + msgType
	if isResponse {
		queue = queue + "_response"
	}
	messages, err := c.channel.Consume(
		queue,
		"",
		c.config.Channel.AutoAck,
		c.config.Channel.Exclusive,
		c.config.Channel.NoLocal,
		c.config.Channel.NoWait,
		nil,
	)
	if err != nil {
		return err
	}
	for msg := range messages{
		msgChannel <- NewMessage(msg.Type, msg.CorrelationId, msg.ReplyTo, msg.Body)
	}
	return nil
}

func (c *Client) CreateQueue(message string, isResponse bool) error {
	queueName := c.name + "." + message
	if isResponse {
		queueName = queueName + "_response"
	}
	_, err := c.channel.QueueDeclare(
		queueName,
		c.config.Queue.Durable,
		c.config.Queue.AutoDel,
		c.config.Queue.Exclusive,
		c.config.Queue.NoWait,
		nil,
	)
	if err != nil {
		return err
	}
	err = c.channel.QueueBind(queueName, "", message, c.config.Queue.NoWait, nil)
	if err != nil {
		return err
	}
	return err
}
