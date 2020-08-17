package internal

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) createAccount(_ context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return nil, errors.New("implement me")
}

func (b *backend) readAccount(_ context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return nil, errors.New("implement me")
}
