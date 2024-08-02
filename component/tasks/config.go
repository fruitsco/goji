package tasks

import "github.com/fruitsco/goji/conf"

type TaskDriver string

const (
	Queue      TaskDriver = "queue"
	CloudTasks TaskDriver = "gcp_cloudtasks"
	NoOp       TaskDriver = "noop"
)

type CloudTasksConfig struct {
	ProjectID string `conf:"project_id"`
	Region    string `conf:"region"`
}

type Config struct {
	Driver TaskDriver `conf:"driver"`

	CloudTasks *CloudTasksConfig `conf:"cloudtasks"`
}

var DefaultConfig = conf.DefaultConfig{
	"tasks.driver": "queue",
}
