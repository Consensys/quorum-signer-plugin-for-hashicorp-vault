package internal

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	util "github.com/jpmorganchase/quorum-go-utils/account"
	"github.com/stretchr/testify/require"
)

func createBackend(t *testing.T) *backend {
	conf := logical.TestBackendConfig()

	b, err := BackendFactory(context.Background(), conf)
	require.NoError(t, err)
	require.NotNil(t, b)

	return b.(*backend)
}

func TestAccountExistenceCheck(t *testing.T) {
	b := createBackend(t)

	storage := &logical.InmemStorage{}
	entry, err := logical.StorageEntryJSON("accounts/myAcct", "96093cadd4bceb60ebdda5b875f5825ef1e91a8e")
	require.NoError(t, err)

	err = storage.Put(context.Background(), entry)
	require.NoError(t, err)

	tests := map[string]struct {
		req  *logical.Request
		want bool
	}{
		"exists": {
			req: &logical.Request{
				Storage: storage,
				Path:    "accounts/myAcct",
			},
			want: true,
		},
		"does not exist": {
			req: &logical.Request{
				Storage: &logical.InmemStorage{},
				Path:    "accounts/does-not-exist",
			},
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := b.accountExistenceCheck(context.Background(), tt.req, nil)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestReadAccount(t *testing.T) {
	b := createBackend(t)

	storage := &logical.InmemStorage{}
	entry, err := logical.StorageEntryJSON("accounts/myAcct", "96093cadd4bceb60ebdda5b875f5825ef1e91a8e")
	require.NoError(t, err)

	err = storage.Put(context.Background(), entry)
	require.NoError(t, err)

	req := &logical.Request{
		Storage: storage,
		Path:    "accounts/myAcct",
	}

	resp, err := b.readAccount(context.Background(), req, nil)
	require.NoError(t, err)
	require.Equal(t, "96093cadd4bceb60ebdda5b875f5825ef1e91a8e", resp.Data["addr"].(string))
}

func TestReadAccount_AccountNotFound(t *testing.T) {
	b := createBackend(t)

	storage := &logical.InmemStorage{}

	req := &logical.Request{
		Storage: storage,
		Path:    "accounts/myAcct",
	}

	resp, err := b.readAccount(context.Background(), req, nil)
	require.NoError(t, err)
	require.Nil(t, resp)
}

func TestUpdateAccount(t *testing.T) {
	b := createBackend(t)

	_, err := b.updateAccount(nil, nil, nil)
	require.Error(t, err)
	require.EqualError(t, updateUnsupportedErr, err.Error())
}

func TestCreateAccount_CreateNew(t *testing.T) {
	b := createBackend(t)

	storage := &logical.InmemStorage{}

	req := &logical.Request{
		Storage: storage,
		Path:    "accounts/myAcct",
	}

	d := &framework.FieldData{
		Raw: map[string]interface{}{
			"acctID": "myAcct",
		},
		Schema: map[string]*framework.FieldSchema{
			"acctID": {
				Type: framework.TypeString,
			},
		},
	}

	resp, err := b.createAccount(context.Background(), req, d)
	require.NoError(t, err)

	// addr in response
	respAddr := resp.Data["addr"].(string)
	require.NotEmpty(t, respAddr)
	addrByt, err := hex.DecodeString(respAddr)
	require.NoError(t, err)
	require.Len(t, addrByt, 20)

	// addr stored in vault
	addrSE, err := storage.Get(context.Background(), "accounts/myAcct")
	require.NoError(t, err)

	var storedAddr string
	err = addrSE.DecodeJSON(&storedAddr)
	require.NoError(t, err)

	require.Equal(t, respAddr, storedAddr)

	// key stored separately in vault
	keySE, err := storage.Get(context.Background(), "keys/myAcct")
	require.NoError(t, err)

	var storedKey string
	err = keySE.DecodeJSON(&storedKey)
	require.NoError(t, err)

	key, err := util.NewKeyFromHexString(storedKey)
	require.NoError(t, err)
	addrFromKey, err := util.PrivateKeyToAddress(key)
	require.NoError(t, err)

	require.Equal(t, respAddr, addrFromKey.ToHexString())
}

func TestCreateAccount_ImportExisting(t *testing.T) {
	toImport := "a0379af19f0b55b0f384f83c95f668ba600b78f487f6414f2d22339273891eec"
	wantAddr := "4d6d744b6da435b5bbdde2526dc20e9a41cb72e5"

	b := createBackend(t)

	storage := &logical.InmemStorage{}

	req := &logical.Request{
		Storage: storage,
		Path:    "accounts/myAcct",
	}

	d := &framework.FieldData{
		Raw: map[string]interface{}{
			"acctID": "myAcct",
			"import": toImport,
		},
		Schema: map[string]*framework.FieldSchema{
			"acctID": {
				Type: framework.TypeString,
			},
			"import": {
				Type: framework.TypeString,
			},
		},
	}

	resp, err := b.createAccount(context.Background(), req, d)
	require.NoError(t, err)

	// addr in response
	respAddr := resp.Data["addr"].(string)
	require.NotEmpty(t, respAddr)
	require.Equal(t, wantAddr, respAddr)

	// addr stored in vault
	addrSE, err := storage.Get(context.Background(), "accounts/myAcct")
	require.NoError(t, err)

	var storedAddr string
	err = addrSE.DecodeJSON(&storedAddr)
	require.NoError(t, err)

	require.Equal(t, respAddr, storedAddr)

	// key stored separately in vault
	keySE, err := storage.Get(context.Background(), "keys/myAcct")
	require.NoError(t, err)

	var storedKey string
	err = keySE.DecodeJSON(&storedKey)
	require.NoError(t, err)
	require.Equal(t, toImport, storedKey)
}

func TestListAccountIDs(t *testing.T) {
	b := createBackend(t)

	storage := &logical.InmemStorage{}
	entry, err := logical.StorageEntryJSON("accounts/myAcct", "96093cadd4bceb60ebdda5b875f5825ef1e91a8e")
	require.NoError(t, err)

	err = storage.Put(context.Background(), entry)
	require.NoError(t, err)

	entry, err = logical.StorageEntryJSON("accounts/anotherAcct", "96093cadd4bceb60ebdda5b875f5825ef1e91a8e")
	require.NoError(t, err)

	err = storage.Put(context.Background(), entry)
	require.NoError(t, err)

	req := &logical.Request{
		Storage: storage,
		Path:    "accounts/",
	}

	resp, err := b.listAccountIDs(context.Background(), req, nil)
	require.NoError(t, err)
	ids := resp.Data["keys"].([]string)
	require.Len(t, ids, 2)
	require.Contains(t, ids, "myAcct")
	require.Contains(t, ids, "anotherAcct")
}

func TestDelete(t *testing.T) {
	want := []byte{21, 228, 169, 48, 162, 94, 71, 55, 85, 214, 104, 193, 92, 14, 27, 132, 111, 18, 108, 11, 194, 150, 169, 254, 177, 54, 67, 10, 14, 208, 100, 250, 123, 166, 26, 0, 44, 215, 237, 186, 32, 198, 241, 77, 206, 214, 249, 124, 212, 36, 249, 4, 171, 87, 68, 147, 238, 96, 8, 180, 122, 172, 175, 38, 1}
	wantHex := hex.EncodeToString(want)
	fmt.Println(wantHex)
}

func TestSign(t *testing.T) {
	b := createBackend(t)

	key := "a0379af19f0b55b0f384f83c95f668ba600b78f487f6414f2d22339273891eec"
	toSign := "bc4c915d69896b198f0292a72373a2bdcd0d52bccbfcec11d9c84c0fff71b0bc"
	wantSig := "f68df2227e39c9ba87baea5966f0c502b038031b10a39e96a721cd270700362d54bae75dcf035a180c17a3a8cf760bfa91a0a41969c0a1630ba6d20e06aa1a8501"

	storage := &logical.InmemStorage{}

	entry, err := logical.StorageEntryJSON("keys/myAcct", key)
	require.NoError(t, err)

	err = storage.Put(context.Background(), entry)
	require.NoError(t, err)

	req := &logical.Request{
		Storage: storage,
		Path:    "sign/myAcct",
	}

	d := &framework.FieldData{
		Raw: map[string]interface{}{
			"acctID": "myAcct",
			"sign":   toSign,
		},
		Schema: map[string]*framework.FieldSchema{
			"acctID": {
				Type: framework.TypeString,
			},
			"sign": {
				Type: framework.TypeString,
			},
		},
	}

	resp, err := b.sign(context.Background(), req, d)
	require.NoError(t, err)
	require.Equal(t, wantSig, resp.Data["sig"])
}
