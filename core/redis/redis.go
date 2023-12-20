package redis

import (
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var (
	ErrConnectionNotConfigured = errors.New("connection not configured")
)

type Connection struct {
	*redis.Client
}

func NewConnection(config *ConnectionConfig) *Connection {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	return &Connection{
		Client: client,
	}
}

type Redis struct {
	config      *Config
	connections map[ConnectionName]*Connection
}

type RedisParams struct {
	fx.In

	Config *Config
}

func New(params RedisParams) *Redis {
	// init default connections
	connections := map[ConnectionName]*Connection{
		"default": NewConnection(params.Config.Connections[params.Config.DefaultConnection]),
	}

	return &Redis{
		config:      params.Config,
		connections: connections,
	}
}

func (r *Redis) resolveConnection(name ConnectionName) (*Connection, error) {
	if conn, ok := r.connections[name]; ok {
		return conn, nil
	}

	if config, ok := r.config.Connections[name]; ok {
		conn := NewConnection(config)
		r.connections[name] = conn
		return conn, nil
	}

	return nil, ErrConnectionNotConfigured
}

func (r *Redis) Default() (*Connection, error) {
	return r.resolveConnection(r.config.DefaultConnection)
}

func (r *Redis) Connection(name ConnectionName) (*Connection, error) {
	return r.resolveConnection(name)
}
