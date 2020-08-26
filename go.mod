module github.com/jpmorganchase/hashicorp-vault-plugin-quorum-signer

go 1.13

replace github.com/jpmorganchase/quorum-go-utils => /Users/chrishounsom/quorum-go-utils

require (
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/hashicorp/vault/sdk v0.1.13
	github.com/jpmorganchase/quorum-go-utils v0.0.0
	github.com/jpmorganchase/quorum/crypto/secp256k1 v0.0.0-20200819121702-fa97b9d9ee78
	github.com/stretchr/testify v1.6.1
)
