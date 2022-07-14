// Package consumer defines the Consumer interface.
package consumer

import "context"

//go:generate mockgen -self_package=iam-pump/internal/pumpserver/consumer -destination mock_consumer.go -package consumer iam-pump/internal/pumpserver/consumer Consumer

// MsgHandler is the callback for handling the received message.
type MsgHandler func(ctx context.Context, message []byte)

var consumer Consumer

// Consumer defines the behavior of a consumer.
type Consumer interface {
	// Start starts the message consuming loop.
	Start(context.Context)
	// Stop stops the message consuming loop.
	Stop(context.Context) error
}

func SetConsumer(c Consumer) {
	consumer = c
}

func GetConsumer() Consumer {
	return consumer
}
