package queue

import "context"

type MessageHandler interface {
	HandleMessage(context.Context, Message) error
}

type MessageHandlerFunc func(context.Context, Message) error

func (f MessageHandlerFunc) HandleMessage(ctx context.Context, message Message) error {
	return f(ctx, message)
}

var _ = MessageHandler(MessageHandlerFunc(nil))
