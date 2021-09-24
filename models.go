package rabbitmq

import "github.com/streadway/amqp"

type (
	MessagePost interface {
		SendMessage(from, to, msg, ID, body string) error
		ReceiveMessage(msg string) (string, string, error)
	}
	Client struct {
		name       string
		config     *config
		connection *amqp.Connection
		channel    *amqp.Channel
	}
	Message struct {
		text   string
		id     string
		body    []byte
		replyTo string
	}
)
