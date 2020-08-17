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
			b.accountPath(),
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

func (b *backend) accountPath() *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("accounts/%s", framework.GenericNameRegex("acctID")),

		Fields: map[string]*framework.FieldSchema{
			"acctID": {
				Type:        framework.TypeString,
				Description: "Specifies the path of the account.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			// TODO(cjh) create vs update? ExistenceCheck have an impact?
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.createAccount,
				Summary:  "Create/update account",
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.readAccount,
				Summary:  "Read account",
			},
		},

		// TODO(cjh) needed?
		//ExistenceCheck: b.handleExistenceCheck,
	}
}
