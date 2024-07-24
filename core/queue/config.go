package queue

import "github.com/fruitsco/goji/conf"

type QueueDriver string

const (
	PubSub QueueDriver = "pubsub"
	NoOp   QueueDriver = "noop"
)

type PubSubConfig struct {
	ProjectID    string  `conf:"project_id"`
	EmulatorHost *string `conf:"emulator_host"`
}

type Config struct {
	Driver QueueDriver   `conf:"driver"`
	PubSub *PubSubConfig `conf:"pubsub"`
}

var DefaultConfig = conf.DefaultConfig{
	"queue.driver": "noop",
}
