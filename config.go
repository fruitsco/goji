package goji

// MARK: - Root Config

type AppConfig struct {
	Environment Environment `conf:"env"`
	Name        string      `conf:"name"`
}

type LogConfig struct {
	Name  string `conf:"name"`
	Level string `conf:"level"`
}

type RootConfig[C any] struct {
	App *AppConfig `conf:"app"`
	Log *LogConfig `conf:"log"`

	Child C `conf:",squash"`
}
