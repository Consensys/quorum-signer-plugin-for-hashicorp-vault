# Quorum Signer plugin for Hashicorp Vault

The Quorum Signer plugin is a [custom plugin backend for Hashicorp Vault](https://www.vaultproject.io/docs/plugin) that adds a new `quorum-signer` secret-engine type to Hashicorp Vault.  

The `quorum-signer` secret-engine creates and stores Quorum accounts (including private keys) and can be used to sign data with those accounts.  

When used in conjunction with the [Hashicorp Vault plugin for Quorum](https://github.com/consensys/quorum-account-plugin-hashicorp-vault), a Quorum (or `clef`) node can sign transactions and data as normal, with the added security benefit that account private keys never leave the boundaries of Vault and never have to be directly managed.

## Building
```shell
make
```

## Quickstart 
> Note: Starting the Vault server in dev mode requires very little setup and is useful for experimentation/testing.  It is insecure and does not persist data between restarts so should not be used for production.

> Note: Using plugins with a non-dev mode Vault server requires additional Vault configuration and for the plugin to be registered before it can be used.  See [Plugin Registration](https://www.vaultproject.io/docs/internals/plugins#plugin-registration) for more info.

```shell
make
```
```shell
vault server -dev -dev-root-token-id=root \
    -dev-plugin-dir=/path/to/hashicorp-vault-signing-plugin/build
``` 

The output should include something similar to the following to indicate the plugin is available:
```shell
The following dev plugins are registered in the catalog:
    - quorum-signer-<VERSION>
```

In another terminal:
```shell
export VAULT_TOKEN=root
vault secrets enable -path quorum-signer quorum-signer-<VERSION>
```

The `quorum-signer` secret-engine will now be available for use. 

### Vault non-dev mode
1. Add `plugin_directory` and `api_addr` fields to `config.hcl`, e.g.: 
    ```
    plugin_directory = "/hashicorp-vault-signing-plugin/build"
    api_addr = "https//localhost:8200"
    ``` 
1. Register the plugin in Vault
    ```shell
    vault write sys/plugins/catalog/secret/quorum-signer-<VERSION> \
        sha256=<BINARY SHA256SUM> \
        command="quorum-signer-<VERSION> --ca-cert=<CA CERT> --client-cert=<CLIENT CERT> --client-key=<CLIENT KEY>"
    ```
   * `<BINARY SHA256SUM>`: Hash of plugin binary (e.g. from `shasum -a 256 /hashicorp-vault-signing-plugin/build/quorum-signer-<VERSION>`)
   * `<CA CERT>`, `<CLIENT CERT>`, `<CLIENT KEY>`: The plugin acts as a client to the Vault server.  If TLS is configured on the Vault server then the paths to the necessary client TLS certs must be provided

## API
The `quorum-signer` secret-engine stores accounts with a user-defined `acctID` (e.g. `myAcct`).  Interacting with accounts is made possible through the plugin's API.

### Create new account
> Note: Overwriting existing secrets (i.e. using the same `acctID` is not supported)

```shell
vault write -f quorum-signer/accounts/<acctID>

Key     Value
---     -----
addr    874f98d93427b145fcf1bb2c34f733f6c14597df 
```

### Import existing account
> Note: Overwriting existing secrets (i.e. using the same `acctID` is not supported)

```shell
vault write quorum-signer/accounts/<acctID> import=1fe8f1ad4053326db20529257ac9401f2e6c769ef1d736b8c2f5aba5f787c72b

Key     Value
---     -----
addr    6038dc01869425004ca0b8370f6c81cf464213b3 
```

* `import`: hex-encoded private key

### Get public account data
```shell
 vault read quorum-signer/accounts/<acctID>

Key     Value
---     -----
addr    874f98d93427b145fcf1bb2c34f733f6c14597df
```

### Sign data with an account
> Note: The `quorum-signer` is a "dumb" signer - it simply signs the provided data with the specified account.  Quorum data is prefixed and hashed before it is signed (e.g. [EIP-191](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-191.md).  Any data sent to the `quorum-signer` for signing should have already been prefixed and hashed.  
>
> This is handled automatically when using `quorum-signer` in conjunction with the [Hashicorp Vault plugin for Quorum](https://github.com/consensys/quorum-account-plugin-hashicorp-vault).

```shell
vault read quorum-signer/sign/myAcct sign=bc4c915d69896b198f0292a72373a2bdcd0d52bccbfcec11d9c84c0fff71b0bc

Key    Value
---    -----
sig    01b4402e23ae8cbff32e708ab485f8e708ccd8b47707b91fad42a5b6353b31ba02579620df93c1a6a189303fcf7a8095eb9c24a7bbc0039ab34e7df7bb6f3b5a01
```

* `sign`: hex-encoded data (prefixed and hashed) to be signed

## Further reading
* [Hashicorp Vault's plugin system](https://www.vaultproject.io/docs/internals/plugins)
