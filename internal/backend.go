package internal

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"strings"
)

const (
	acctPath = "accounts"
	signPath = "sign"
	keyPath  = "keys"
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
			b.signPath(),
		},
		PathsSpecial: &logical.Paths{
			// paths to encrypt when sealed
			SealWrapStorage: []string{
				fmt.Sprintf("%s/", acctPath),
				fmt.Sprintf("%s/", keyPath),
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
		Pattern: fmt.Sprintf("%s/?$", acctPath),

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
		Pattern: fmt.Sprintf("%s/%s", acctPath, framework.GenericNameRegex("acctID")),

		Fields: map[string]*framework.FieldSchema{
			"acctID": {
				Type:        framework.TypeString,
				Description: "Specifies the path of the account.",
			},
			"import": {
				Type:        framework.TypeString,
				Description: "(optional) A hex-encoded private key to imported and store at the specified path.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.readAccount,
				Summary:  "Read account address",
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.createAccount,
				Summary:  "Generate and store new Quorum account, or import existing account by using the 'import' field.",
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.updateAccount,
				Summary:  "Generate and store new Quorum account, or import existing account by using the 'import' field.",
			},
		},
		ExistenceCheck: b.accountExistenceCheck, // determines whether create or update operation is called
	}
}

func (b *backend) signPath() *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("%s/%s", signPath, framework.GenericNameRegex("acctID")),

		Fields: map[string]*framework.FieldSchema{
			"acctID": {
				Type:        framework.TypeString,
				Description: "Specifies the path of the account.",
			},
			"sign": {
				Type:        framework.TypeString,
				Description: "Hex-encoded payload to be signed.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.sign,
				Summary:  "Sign data with account, returns hex-encoded signature in r,s,v format where v is 0 or 1",
			},
		},
	}
}
