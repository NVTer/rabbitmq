package RabbitMQ

func NewMessage(text, ID string, body []byte) *Message {
	return &Message{text: text, id: ID, body: body}
}

func (m *Message) GetText() string {
	return m.text
}

func (m *Message) GetBody() []byte {
	return m.body
}

func (m *Message) GetID() string {
	return m.id
}
