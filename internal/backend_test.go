package internal

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestPathRoutingAndSupportedOperations(t *testing.T) {
	b, err := BackendFactory(context.Background(), &logical.BackendConfig{})
	require.NoError(t, err)

	storage := &logical.InmemStorage{}

	tests := []struct {
		name       string
		path       string
		operations []logical.Operation
	}{
		{
			name:       "rootpath rw",
			path:       "accounts/mySecret",
			operations: []logical.Operation{logical.ReadOperation, logical.CreateOperation, logical.UpdateOperation},
		},
		{
			name:       "rootpath list",
			path:       "accounts/",
			operations: []logical.Operation{logical.ListOperation},
		},
		{
			name:       "subpath rw",
			path:       "accounts/myApp/appSecret",
			operations: []logical.Operation{logical.ReadOperation, logical.CreateOperation, logical.UpdateOperation},
		},
		{
			name:       "subpath list",
			path:       "accounts/myApp/",
			operations: []logical.Operation{logical.ListOperation},
		},
		{
			name:       "rootpath sign",
			path:       "sign/mySecret",
			operations: []logical.Operation{logical.ReadOperation},
		},
		{
			name:       "subpath sign",
			path:       "sign/myApp/appSecret",
			operations: []logical.Operation{logical.ReadOperation},
		},
	}

	for _, tt := range tests {
		for _, op := range tt.operations {
			testName := fmt.Sprintf("%s_%s", tt.name, op)
			t.Run(testName, func(t *testing.T) {
				req := &logical.Request{
					Storage:   storage,
					Path:      tt.path,
					Operation: op,
				}

				// sign/ requires request data
				if strings.HasPrefix(tt.path, "sign/") {
					req.Data = map[string]interface{}{
						"sign": "7d15728d30727d67a3257e6bbd4724c4d31f830f017fd0e0d2d802c14bdf408d",
					}
				}

				_, err := b.HandleRequest(context.Background(), req)

				if logical.UpdateOperation == op {
					// updating is unsupported.
					// we still want to make sure requests are routed correctly and our error is returned.
					require.ErrorIs(t, err, updateUnsupportedErr)
				} else {
					require.NoError(t, err)
				}
			})
		}
	}
}
