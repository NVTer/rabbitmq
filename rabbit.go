package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
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

func (c *Client) SendMessage(msg, corrID string, body []byte) error {
	replyTo := c.name + "." + msg + "_response"
	return c.channel.Publish(msg, "", c.config.Mandatory, c.config.Immediate, amqp.Publishing{
		ContentType:   "application/json",
		Type:          msg,
		CorrelationId: corrID,
		ReplyTo:       replyTo,
		Body:          body,
	})
}

func (c *Client) SendReply(consumer, msg, corrID string, body []byte) error {
	return c.channel.Publish("", consumer, c.config.Mandatory, c.config.Immediate, amqp.Publishing{
		ContentType:   "application/json",
		Type:          msg,
		CorrelationId: corrID,
		Body:          body,
	})
}

func (c *Client) StartConsumer(msg string, msgChannel chan *Message, isResponse bool) error {
	queue := c.name + "." + msg
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
	for msg := range messages {
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
	if !isResponse {
		err = c.channel.QueueBind(queueName, "", message, c.config.Queue.NoWait, nil)
	}
	if err != nil {
		return err
	}
	return err
}

func (c *Client) CreateExchange(name string) error {
	return c.channel.ExchangeDeclare(
		name,
		c.config.Exchange.Type,
		c.config.Exchange.Durable,
		c.config.Exchange.AutoDel,
		c.config.Exchange.Internal,
		c.config.Exchange.NoWait,
		nil,
	)
}

func (c *Client) Ping() error {
	notify := c.channel.NotifyClose(make(chan *amqp.Error))
	select {
	case err := <-notify:
		return err
	case <-time.After(5 * time.Second):
		return nil
	}
}
