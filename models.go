package RabbitMQ

import "github.com/streadway/amqp"

type (
	MessagePost interface {
		SendMessage(from, to, msg, ID, body string) error
		ReceiveMessage(msg string) (string, string, error)
	}

	Client struct {
		name    string
		channel *amqp.Channel
	}
	Message struct {
		text string
		id   string
		body []byte
	}
)
