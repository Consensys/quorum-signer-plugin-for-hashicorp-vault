package internal

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(context.Context, *logical.BackendConfig) (logical.Backend, error) {
	return nil, errors.New("implement me")
}
