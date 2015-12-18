package server

type PubSub interface {
	Subscribe(topic string) (chan []byte, error)
	Unsubscribe(topic string, ch chan []byte) error
	Publish(topic string, msg string) error
}
