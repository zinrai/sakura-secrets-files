# sakura-secrets-files

A command-line tool that makes secrets from [SAKURA Cloud Secret Manager](https://cloud.sakura.ad.jp/products/secrets-manager/) present or absent at declared file paths.

## Features

- A YAML manifest declares which secret goes to which path
- `present` compares each secret against the current file content and reports the decision per file: `create`, `update` or `unchanged`. Files are written atomically with 0600 permissions, and only when the content differs
- `absent` removes the declared paths, so a wrapper can bound the lifetime of materialized secrets to a single operation
- `absent` needs no credentials and no network. Removing files is a purely local operation
- `-dry-run` reports the same decisions without changing anything
- Fail-fast with no retry logic. Errors exit immediately with a detailed message

## Requirements

- SAKURA Cloud account with Secret Manager access
- Valid API credentials (static API keys or a service principal)

## Configuration

### Environment Variables

API credentials are resolved by [sacloud-sdk-go](https://github.com/sacloud/sacloud-sdk-go). Set either static API keys:

```bash
$ export SAKURA_ACCESS_TOKEN="your-access-token"
$ export SAKURA_ACCESS_TOKEN_SECRET="your-access-token-secret"
```

or service principal credentials:

```bash
$ export SAKURA_SERVICE_PRINCIPAL_ID="your-service-principal-id"
$ export SAKURA_SERVICE_PRINCIPAL_KEY_ID="your-key-id"
$ export SAKURA_PRIVATE_KEY_PATH="/path/to/private-key.pem"
```

### Configuration File

Create a YAML manifest (e.g., `secrets.yaml`):

```yaml
vault:
  id: "123456789012"  # Your vault resource ID
  zone: is1a          # Optional (default: is1a)

secrets:
  - name: production-db-password
    dest: roles/database/files/db_password

  - name: production-api-key
    dest: roles/application/files/api_key
```

## Usage

Write the declared secrets to their paths:

```bash
$ sakura-secrets-files present -config secrets.yaml
[create] production-db-password -> roles/database/files/db_password
[update] production-api-key -> roles/application/files/api_key
```

Remove them:

```bash
$ sakura-secrets-files absent -config secrets.yaml
[removed] roles/database/files/db_password
[removed] roles/application/files/api_key
```

### Options

Both subcommands take:

- `-config <path>`: Path to the manifest (required)
- `-dry-run`: Report decisions without changing anything

## Error Handling

- If any secret fails to fetch or write, the tool exits immediately with exit code 1
- `absent` treats an already missing path as success and reports it as `[absent]`
- No retry logic. Re-run the command if it fails

## Security

- Secrets are written with 0600 permissions (owner read/write only)
- Atomic writes using temporary files ensure no partial writes
- API credentials are never logged or written to disk

## License

This project is licensed under the [MIT License](./LICENSE).
