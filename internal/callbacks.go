package internal

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/jpmorganchase/quorum/crypto/secp256k1"
)

type hexAccountData struct {
	HexAddress string
	HexKey     string
}

func (b *backend) readAccount(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("reading account", "path", req.Path)

	// TODO(cjh) perhaps we should store the addr and key separately so that we only need to Get the addr - this might
	//  get complicated when we consider versioning.  Perhaps store ID -> addr, and addr -> key to get around this?
	storageEntry, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return nil, err
	}

	hexAccountData := new(hexAccountData)
	if err := storageEntry.DecodeJSON(hexAccountData); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"addr": hexAccountData.HexAddress,
		},
	}, nil
}

func (b *backend) createAccount(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var (
		hexAccountData *hexAccountData
		err            error
	)

	if rawKey, ok := d.GetOk("import"); ok {
		b.Logger().Info("importing existing account", "path", req.Path)
		rawKeyStr, ok := rawKey.(string)
		if !ok {
			return nil, errors.New("key to import must be a valid string")
		}
		hexAccountData, err = rawKeyToHexAccountData(rawKeyStr)
		if err != nil {
			return nil, fmt.Errorf("unable to import account: %v", err)
		}
	} else {
		b.Logger().Info("creating new account", "path", req.Path)

		hexAccountData, err = generateAccount()
		if err != nil {
			return nil, fmt.Errorf("unable to generate new account: %v", err)
		}
	}

	storageEntry, err := logical.StorageEntryJSON(req.Path, hexAccountData)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, storageEntry); err != nil {
		return nil, fmt.Errorf("unable to store account: %v", err)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
	}
	resp.Data["addr"] = hexAccountData.HexAddress

	return resp, nil
}

func generateAccount() (*hexAccountData, error) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	defer zeroKey(key)

	return keyToHexAccountData(key)
}

func rawKeyToHexAccountData(rawKey string) (*hexAccountData, error) {
	key, err := NewKeyFromHexString(rawKey)
	if err != nil {
		return nil, err
	}
	defer zeroKey(key)

	return keyToHexAccountData(key)
}

func keyToHexAccountData(key *ecdsa.PrivateKey) (*hexAccountData, error) {
	addr, err := PrivateKeyToAddress(key)
	if err != nil {
		return nil, err
	}

	hexKey, err := PrivateKeyToHexString(key)
	if err != nil {
		return nil, err
	}

	return &hexAccountData{
		HexAddress: addr.ToHexString(),
		HexKey:     hexKey,
	}, nil
}

func (b *backend) listAccountIDs(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("listing account IDs", "path", req.Path)

	ids, err := req.Storage.List(ctx, req.Path)
	if err != nil {
		return nil, err
	}

	b.Logger().Info("account IDs retrieved from storage", "IDs", ids)

	return logical.ListResponse(ids), nil
}
