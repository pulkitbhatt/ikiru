package publisher

import "context"

type Message struct {
	Destination string
	Payload     []byte
	Headers     map[string]string
	Key         string
}

type Publisher interface {
	Publish(ctx context.Context, msg Message) error
}
