package google

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type Credentials struct {
	credentials *google.Credentials
}

func NewCredentials(ctx context.Context, scopes []string) *Credentials {
	credentials, err := google.FindDefaultCredentials(ctx, scopes...)

	if err != nil {
		fmt.Println(err)
	}

	return &Credentials{
		credentials: credentials,
	}
}

func (c *Credentials) ClientOption() option.ClientOption {
	var clientOption option.ClientOption

	if c.credentials != nil {
		return option.WithCredentials(c.credentials)
	}

	return clientOption
}
