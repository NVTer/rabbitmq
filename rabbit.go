package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

func NewClient(serviceName string, config *Config) (*Client, error) {
	connect, err := amqp.Dial(NewRabbitMQConn(config))
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

func NewRabbitMQConn(cfg *Config) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
}

func (c *Client) SendMessage(msgType, corrID string, body []byte) error {
	replyTo := c.name + "." + msgType + "_response"
	return c.channel.Publish(msgType, "", c.config.Mandatory, c.config.Immediate, amqp.Publishing{
		ContentType:   "application/json",
		Type:          msgType,
		CorrelationId: corrID,
		ReplyTo:       replyTo,
		Body:          body,
	})
}

func (c *Client) SendReply(consumer, msgType, corrID string, body []byte) error {
	return c.channel.Publish("", consumer, c.config.Mandatory, c.config.Immediate, amqp.Publishing{
		ContentType:   "application/json",
		Type:          msgType,
		CorrelationId: corrID,
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


