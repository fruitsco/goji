package redis

import "github.com/fruitsco/goji/conf"

type ConnectionName string

const (
	DefaultConnectionName ConnectionName = "default"
)

type ConnectionConfig struct {
	Name     string `conf:"name"`
	Host     string `conf:"host"`
	Port     int    `conf:"port"`
	Password string `conf:"password"`
	DB       int    `conf:"db"`
}

type Config struct {
	DefaultConnection ConnectionName `conf:"connection"`

	Connections map[ConnectionName]*ConnectionConfig `conf:"connections"`
}

var DefaultConfig = conf.DefaultConfig{
	"redis.connection":                   "default",
	"redis.connections.default.name":     "default",
	"redis.connections.default.host":     "localhost",
	"redis.connections.default.port":     "6379",
	"redis.connections.default.password": "",
	"redis.connections.default.db":       "0",
}
