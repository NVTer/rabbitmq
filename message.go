package rabbitmq

func NewMessage(text, ID, sender string, body []byte) *Message {
	return &Message{text: text, id: ID, body: body, replyTo: sender}
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

func (m *Message) GetReplyTo() string {
	return m.replyTo
}
