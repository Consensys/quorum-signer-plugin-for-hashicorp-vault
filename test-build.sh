#!/bin/zsh
set -x
# rebuild the plugin, clean up any existing vault server, start a new vault server

value=$(cat vlt.pid)
kill ${value}

go build -o bin/quorum-signer main.go

vault server -dev -dev-root-token-id=root -dev-plugin-dir=bin &

echo $! > vlt.pid

export VAULT_ADDR="http://127.0.0.1:8200"

vault login root

vault secrets enable quorum-signer

vault write -force quorum-signer/accounts/myAcct
vault write -force quorum-signer/accounts/myAcct
vault write -force quorum-signer/accounts/myOtherAcct
#vault write quorum-signer/accounts/myImportedAcct rawKey=1fe8f1ad4053326db20529257ac9401f2e6c769ef1d736b8c2f5aba5f787c72b

vault read quorum-signer/accounts/myAcct
#vault read quorum-signer/accounts/myOtherAcct
#vault read quorum-signer/accounts/myImportedAcct

vault list quorum-signer/accounts
vault list quorum-signer/accounts/
#vault list quorum-signer/accounts/name
#vault list quorum-signer/accounts/address

set +x
