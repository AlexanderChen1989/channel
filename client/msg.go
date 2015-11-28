package client

type Msg struct {
	Topic string `json:"topic"`
	Event string `json:"event"`
	Ref   string `json:"ref"`

	Payload interface{} `json:"payload"`
}
