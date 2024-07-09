package crypt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"go.uber.org/fx"
)

// CryptoParams is a struct that holds the dependencies of the Crypto module.
type CryptoParams struct {
	fx.In

	// Config is the configuration for the Crypto module.
	Config *Config

	// KeyProvider is the key provider for the Crypto module.
	KeyProvider KeyProvider
}

// Crypto is a struct that provides encryption and decryption functionality.
type Crypto struct {
	config      *Config
	keyProvider KeyProvider
}

// New creates a new instance of the Crypto module.
func New(params CryptoParams) (*Crypto, error) {
	return &Crypto{
		config:      params.Config,
		keyProvider: params.KeyProvider,
	}, nil
}

// Encrypt encrypts the given data using the key with the given name.
func (c *Crypto) Encrypt(ctx context.Context, data []byte, keyName string) (Capsule, error) {
	key, err := c.keyProvider.GetKey(ctx, keyName)
	if err != nil {
		return Capsule{}, fmt.Errorf("failed to get encryption key: %w", err)
	}

	aes, err := aes.NewCipher([]byte(key.Data))
	if err != nil {
		return Capsule{}, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return Capsule{}, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return Capsule{}, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return Capsule{
		Data:       ciphertext,
		KeyName:    key.Name,
		KeyVersion: key.Version,
	}, nil
}

// Decrypt decrypts the given capsule.
func (c *Crypto) Decrypt(ctx context.Context, capsule Capsule) ([]byte, error) {
	key, err := c.keyProvider.GetKeyVersion(ctx, capsule.KeyName, capsule.KeyVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption key: %w", err)
	}

	aes, err := aes.NewCipher([]byte(key.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := capsule.Data[:nonceSize], capsule.Data[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// Recrypt re-encrypts the given capsule with the latest version of the key.
func (c *Crypto) Recrypt(ctx context.Context, capsule Capsule) (Capsule, error) {
	plaintext, err := c.Decrypt(ctx, capsule)
	if err != nil {
		return Capsule{}, fmt.Errorf("failed to decrypt capsule: %w", err)
	}

	return c.Encrypt(ctx, plaintext, capsule.KeyName)
}
