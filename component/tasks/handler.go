package tasks

import "context"

type Handler interface {
	HandleTask(context.Context, *Task) error
}

type HandlerFunc func(context.Context, *Task) error

func (f HandlerFunc) HandleTask(ctx context.Context, message *Task) error {
	return f(ctx, message)
}

var _ = Handler(HandlerFunc(nil))
