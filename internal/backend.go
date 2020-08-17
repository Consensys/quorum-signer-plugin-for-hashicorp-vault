package internal

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"strings"
)

type backend struct {
	*framework.Backend
}

func BackendFactory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	var b backend
	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace("Creates and stores Quorum accounts.  Signs data using those accounts.\n"),
		BackendType: logical.TypeLogical,
		Paths: []*framework.Path{
			b.accountsPath(),
			b.accountIDPath(),
		},
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"accounts/", // paths to encrypt when sealed
			},
		},
	}

	if err := b.Backend.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *backend) accountsPath() *framework.Path {
	return &framework.Path{
		Pattern: "accounts/?$",

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.listAccountIDs,
				Summary:  "List account IDs",
			},
		},
	}
}

func (b *backend) accountIDPath() *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("accounts/%s", framework.GenericNameRegex("acctID")),

		Fields: map[string]*framework.FieldSchema{
			"acctID": {
				Type:        framework.TypeString,
				Description: "Specifies the path of the account.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.createAccount,
				Summary:  "Create/update account",
			},
		},
	}
}
