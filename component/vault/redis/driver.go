package vaultredis

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/redis"
	"github.com/fruitsco/goji/component/vault"
	"github.com/fruitsco/goji/x/driver"
)

// RedisDriver is the driver for Redis
type RedisDriver struct {
	config *vault.RedisConfig
	redis  *redis.Client
	log    *zap.Logger
}

// RedisDriverParams is the parameters for the Redis driver
type RedisDriverParams struct {
	fx.In

	// Config is the configuration for the Redis driver
	Config *vault.RedisConfig

	// Redis is the Redis connection
	Redis *redis.Redis

	// Log is the logger for the Redis driver
	Log *zap.Logger
}

// NewRedisDriverFactory creates a new Redis driver factory
func NewRedisDriverFactory(params RedisDriverParams) driver.FactoryResult[vault.DriverName, vault.Driver] {
	return driver.NewFactory(vault.Redis, func() (vault.Driver, error) {
		return NewRedisDriver(params)
	})
}

// NewRedisDriver creates a new Redis driver
func NewRedisDriver(params RedisDriverParams) (vault.Driver, error) {
	if params.Config == nil {
		return nil, fmt.Errorf("config is required for Redis driver")
	}

	if params.Config.EncryptionKey == "" {
		return nil, fmt.Errorf("encryption key is required for Redis driver")
	}

	if params.Config.ConnectionName == "" {
		params.Config.ConnectionName = redis.DefaultConnectionName
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

var _ = vault.Driver(&RedisDriver{})

// CreateSecret creates a new secret in Redis
func (d *RedisDriver) CreateSecret(
	ctx context.Context,
	name string,
	payload []byte,
) (vault.Secret, error) {
	encryptedPayload, err := d.encrypt(payload)
	if err != nil {
		return vault.Secret{}, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	res := d.redis.LPush(ctx, d.getKeyName(name), encryptedPayload)
	if res.Err() != nil {
		return vault.Secret{}, res.Err()
	}

	// `lpush` returns the length of the list after the push,
	// which corresponds to the 1-based version number
	version := int(res.Val())

	return vault.Secret{
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
) (vault.Secret, error) {
	encryptedPayload, err := d.encrypt(payload)
	if err != nil {
		return vault.Secret{}, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	res := d.redis.LPush(ctx, d.getKeyName(name), encryptedPayload)

	// `lpush` returns the length of the list after the push,
	// which corresponds to the 1-based version number
	version := int(res.Val())

	return vault.Secret{
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
) (vault.Secret, error) {
	res := d.redis.LIndex(ctx, d.getKeyName(name), int64(-version))
	if res.Err() != nil {
		return vault.Secret{}, res.Err()
	}

	decryptedPayload, err := d.decrypt(res.Val())
	if err != nil {
		return vault.Secret{}, fmt.Errorf("failed to decrypt payload: %w", err)
	}

	return vault.Secret{
		Name:    name,
		Version: version,
		Payload: decryptedPayload,
	}, nil
}

// GetLatestVersion retrieves the latest version of a secret from Redis
func (d *RedisDriver) GetLatestVersion(
	ctx context.Context,
	name string,
) (vault.Secret, error) {
	itemRes := d.redis.LIndex(ctx, d.getKeyName(name), 0)
	if itemRes.Err() != nil {
		return vault.Secret{}, itemRes.Err()
	}

	// decrypt the payload
	decryptedPayload, err := d.decrypt(itemRes.Val())
	if err != nil {
		return vault.Secret{}, fmt.Errorf("failed to decrypt payload: %w", err)
	}

	// get the length of the list, which corresponds to the version number
	lenRes := d.redis.LLen(ctx, d.getKeyName(name))
	if lenRes.Err() != nil {
		return vault.Secret{}, lenRes.Err()
	}

	return vault.Secret{
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

func (d *RedisDriver) encrypt(data []byte) (string, error) {
	aes, err := aes.NewCipher([]byte(d.config.EncryptionKey))
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	encryptedPayload := gcm.Seal(nonce, nonce, data, nil)

	// encode the encrypted payload to base64
	encodedPayload := base64.StdEncoding.EncodeToString(encryptedPayload)

	return encodedPayload, nil
}

func (d *RedisDriver) decrypt(data string) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to base64decode data: %w", err)
	}

	aes, err := aes.NewCipher([]byte(d.config.EncryptionKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := decodedData[:nonceSize], decodedData[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}
