package vault

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"github.com/fruitsco/goji/core/redis"
	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RedisDriverConfig is the configuration for the Redis driver
type RedisDriverConfig struct {
	// ConnectionName is the name of the Redis connection to use
	ConnectionName string `conf:"connection_name"`

	// EncryptionKey is the key to use for encryption
	EncryptionKey string `conf:"encryption_key"`
}

// RedisDriver is the driver for Redis
type RedisDriver struct {
	config *RedisDriverConfig
	redis  *redis.Connection
	log    *zap.Logger
}

// RedisDriverParams is the parameters for the Redis driver
type RedisDriverParams struct {
	fx.In

	// Config is the configuration for the Redis driver
	Config *RedisDriverConfig

	// Redis is the Redis connection
	Redis *redis.Redis

	// Log is the logger for the Redis driver
	Log *zap.Logger
}

// NewRedisDriverFactory creates a new Redis driver factory
func NewRedisDriverFactory(params RedisDriverParams) driver.FactoryResult[DriverName, Driver] {
	return driver.NewFactory(Redis, func() (Driver, error) {
		return NewRedisDriver(params)
	})
}

// NewRedisDriver creates a new Redis driver
func NewRedisDriver(params RedisDriverParams) (Driver, error) {
	if params.Config == nil {
		return nil, fmt.Errorf("config is required for Redis driver")
	}

	if params.Config.EncryptionKey == "" {
		return nil, fmt.Errorf("encryption key is required for Redis driver")
	}

	if params.Config.ConnectionName == "" {
		params.Config.ConnectionName = "default"
	}

	connectionName := redis.ConnectionName(params.Config.ConnectionName)

	connection, err := params.Redis.Connection(connectionName)
	if err != nil {
		return nil, err
	}

	return &RedisDriver{
		redis:  connection,
		config: params.Config,
		log:    params.Log.Named("redis"),
	}, nil
}

var _ = Driver(&RedisDriver{})

// CreateSecret creates a new secret in Redis
func (d *RedisDriver) CreateSecret(
	ctx context.Context,
	name string,
	payload []byte,
) (Secret, error) {
	encryptedPayload, err := d.encrypt(payload)
	if err != nil {
		return Secret{}, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	res := d.redis.LPush(ctx, d.getKeyName(name), string(encryptedPayload))
	if res.Err() != nil {
		return Secret{}, res.Err()
	}

	// `lpush` returns the length of the list after the push,
	// which corresponds to the 1-based version number
	version := int(res.Val())

	return Secret{
		Name:    name,
		Version: version,
		Payload: payload,
	}, nil
}

// AddVersion adds a new version of a secret to Redis
func (d *RedisDriver) AddVersion(
	ctx context.Context,
	name string,
	payload []byte,
) (Secret, error) {
	encryptedPayload, err := d.encrypt(payload)
	if err != nil {
		return Secret{}, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	res := d.redis.LPush(ctx, d.getKeyName(name), string(encryptedPayload))

	// `lpush` returns the length of the list after the push,
	// which corresponds to the 1-based version number
	version := int(res.Val())

	return Secret{
		Name:    name,
		Version: version,
		Payload: payload,
	}, nil
}

// GetVersion retrieves a specific version of a secret from Redis
func (d *RedisDriver) GetVersion(
	ctx context.Context,
	name string,
	version int,
) (Secret, error) {
	res := d.redis.LIndex(ctx, d.getKeyName(name), int64(-version))
	if res.Err() != nil {
		return Secret{}, res.Err()
	}

	decryptedPayload, err := d.decrypt([]byte(res.Val()))
	if err != nil {
		return Secret{}, fmt.Errorf("failed to decrypt payload: %w", err)
	}

	return Secret{
		Name:    name,
		Version: version,
		Payload: decryptedPayload,
	}, nil
}

// GetLatestVersion retrieves the latest version of a secret from Redis
func (d *RedisDriver) GetLatestVersion(
	ctx context.Context,
	name string,
) (Secret, error) {
	itemRes := d.redis.LIndex(ctx, d.getKeyName(name), 0)
	if itemRes.Err() != nil {
		return Secret{}, itemRes.Err()
	}

	decryptedPayload, err := d.decrypt([]byte(itemRes.Val()))
	if err != nil {
		return Secret{}, fmt.Errorf("failed to decrypt payload: %w", err)
	}

	lenRes := d.redis.LLen(ctx, d.getKeyName(name))
	if lenRes.Err() != nil {
		return Secret{}, lenRes.Err()
	}

	return Secret{
		Name:    name,
		Version: int(lenRes.Val()),
		Payload: decryptedPayload,
	}, nil
}

// DeleteSecret deletes a secret from Redis
func (d *RedisDriver) DeleteSecret(ctx context.Context, name string) error {
	res := d.redis.Del(ctx, d.getKeyName(name))
	return res.Err()
}

// getKeyname returns the key name for the secret
func (d *RedisDriver) getKeyName(name string) string {
	return fmt.Sprintf("vault:%s", name)
}

func (d *RedisDriver) encrypt(data []byte) ([]byte, error) {
	aes, err := aes.NewCipher([]byte(d.config.EncryptionKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (d *RedisDriver) decrypt(data []byte) ([]byte, error) {
	aes, err := aes.NewCipher([]byte(d.config.EncryptionKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}
