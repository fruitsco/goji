package notification

type NotificationDriver string

type Topics string

const (
	DomainUpload Topics = "domainUpload"
)

const (
	Slack NotificationDriver = "slack"
	NoOp  NotificationDriver = "noop"
)

type SlackConfig struct {
	URL string `conf:"url"`
}

type Config struct {
	Driver NotificationDriver `conf:"driver"`

	Slack *SlackConfig `conf:"slack"`
}

var DefaultConfig = map[string]any{
	"notification.driver": "slack",
}
