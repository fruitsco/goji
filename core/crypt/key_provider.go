package crypt

import "context"

// Key is a struct that represents an encryption key.
type Key struct {
	// Name is the name of the key.
	Name string

	// Version is the version of the key.
	Version int

	// Data is the key data.
	Data []byte
}

// KeyProvider is an interface for fetching encryption keys.
type KeyProvider interface {
	// GetKey returns the latest version of the key with the given name.
	GetKey(context.Context, string) (Key, error)

	// GetKeyVersion returns the key with the given name and version.
	GetKeyVersion(context.Context, string, int) (Key, error)
}
