package queue

import "context"

type Handler interface {
	HandleMessage(context.Context, Message) error
}

type HandlerFunc func(context.Context, Message) error

func (f HandlerFunc) HandleMessage(ctx context.Context, message Message) error {
	return f(ctx, message)
}

var _ = Handler(HandlerFunc(nil))
