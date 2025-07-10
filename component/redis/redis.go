package redis

import (
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type Client = redis.Client

var (
	ErrConnectionNotConfigured = errors.New("connection not configured")
)

func NewConnection(config *ConnectionConfig) *Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
}

type Redis struct {
	config      *Config
	connections map[ConnectionName]*Client
}

type RedisParams struct {
	fx.In

	Config *Config
}

func New(params RedisParams) *Redis {
	// init default connections
	connections := map[ConnectionName]*Client{
		(DefaultConnectionName): NewConnection(params.Config.Connections[params.Config.DefaultConnection]),
	}

	return &Redis{
		config:      params.Config,
		connections: connections,
	}
}

func (r *Redis) resolveConnection(name ConnectionName) (*Client, error) {
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

func (r *Redis) Default() (*Client, error) {
	return r.resolveConnection(r.config.DefaultConnection)
}

func (r *Redis) Connection(name ConnectionName) (*Client, error) {
	return r.resolveConnection(name)
}
