package beacon

type Message struct {
	msg       string
	frequency int
}

func NewMessage(message string, freq int) *Message {
	return &Message{
		msg:       message,
		frequency: freq,
	}
}
