package queue

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

var DefaultConfig = map[string]any{
	"queue.driver": "noop",
}
