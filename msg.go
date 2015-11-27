package channel

type Message struct {
	Topic   string      `json:"topic"`
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
	Ref     int         `json:"ref"`
}

func Msg(event string, payload interface{}) *Message {
	return nil
}
