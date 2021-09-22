package RabbitMQ

func NewMessage(text, ID string, body []byte) *Message {
	return &Message{text: text,id: ID,body: body}
}
