package rabbitmq

import (
	"github.com/streadway/amqp"
)

func NewClient(serviceName, configFilePath, URL string) (*Client, error) {
	config, err := newConfig(configFilePath)
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

func (c *Client) Send(ID string, body []byte) error {
	return c.channel.Publish("", "", c.config.Mandatory, c.config.Immediate, amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: ID,
		Body:          body,
	})
}

func (c *Client) SendWithReply(msg, ID string, body []byte) (string, error) {
	msgResponse :=msg +"_response"
	replyTo := c.name + "." + msgResponse
	err := c.CreateQueue(msgResponse)
	if err != nil {
		return "", err
	}
	return replyTo, c.channel.Publish(msg, "", c.config.Mandatory, c.config.Immediate, amqp.Publishing{
		ContentType:   "application/json",
		Type:          msg,
		CorrelationId: ID,
		ReplyTo:       replyTo,
		Body:          body,
	})
}

func (c *Client) ReceiveMessage(queue string) (chan *Message, error) {
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
		return nil, err
	}
	m := make(chan *Message)
	for msg := range messages {
		 m<-NewMessage(msg.Type, msg.CorrelationId, msg.ReplyTo, msg.Body)
		break
	}
	return m, nil
}

func (c *Client) receive(){
	return
}

func (c *Client) CreateQueue(message string) error {
	queueName := c.name + "." + message
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
