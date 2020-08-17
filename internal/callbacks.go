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
	hexAddress string
	hexKey     string
}

func (b *backend) createAccount(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("creating account at path %v", req.Path)

	acctID, _, err := d.GetOkErr("acctID")
	if err != nil {
		return nil, err
	}

	// TODO(cjh) check is not entered & remove
	if want := fmt.Sprintf("accounts/%v", acctID); req.Path != want {
		msg := fmt.Sprintf("CHRISSY req.Path was not expected value: want %v, got %v", want, req.Path)
		b.Logger().Error(msg)
		return nil, errors.New(msg)
	}

	hexAccountData, err := generateAccountAsHex()
	if err != nil {
		return nil, fmt.Errorf("unable to generate new account: %v", err)
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
	resp.Data["addr"] = hexAccountData.hexAddress

	return resp, nil
}

func generateAccountAsHex() (*hexAccountData, error) {
	key, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	defer zeroKey(key)

	addr, err := PrivateKeyToAddress(key)
	if err != nil {
		return nil, err
	}

	hexKey, err := PrivateKeyToHexString(key)
	if err != nil {
		return nil, err
	}

	return &hexAccountData{
		hexAddress: addr.ToHexString(),
		hexKey:     hexKey,
	}, nil
}

func (b *backend) listAccountIDs(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	ids, err := req.Storage.List(ctx, req.Path)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
	}
	resp.Data["IDs"] = ids

	return resp, nil
}
