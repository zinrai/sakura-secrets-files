# sakura-secrets-files

A command-line tool that materializes secrets from [SAKURA Cloud Secret Manager](https://cloud.sakura.ad.jp/products/secrets-manager/) as files below a given directory, as declared in a YAML manifest.

## Usage

```bash
$ sakura-secrets-files -config secrets.yaml -base-dir /tmp/tmp.abc123
[write] production-db-password -> /tmp/tmp.abc123/database/db_password
[write] production-api-key -> /tmp/tmp.abc123/application/api_key
```

One line per file on stderr. Every secret in the manifest is fetched and written on every run.

### Options

| Option | Required | Description |
|--------|----------|-------------|
| `-config <path>` | yes | Path to the manifest |
| `-base-dir <dir>` | yes | Directory the `dest` paths are resolved against. It must already exist |
| `-version` | no | Print version and exit |

## Manifest

See [config.yaml.example](./config.yaml.example) for one to copy.

`dest` is a path relative to `-base-dir`. An absolute path, a path that escapes `-base-dir`, and an unknown key anywhere in the manifest are all errors. The latest version of each secret is fetched.

## Credentials

Resolved by [sacloud-sdk-go](https://github.com/sacloud/sacloud-sdk-go) from the environment, and never logged or written to disk. Set either static API keys:

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

## What it writes

- Files with 0600 permissions, into missing intermediate directories created with 0700
- Atomically, through a temporary file in the destination directory
- Nothing outside `-base-dir`. The check is lexical, so it assumes `-base-dir` is a directory you control, and not one another user can plant symlinks in
- Nothing at all until the manifest, every `dest` and `-base-dir` have been checked, so a bad manifest fails without a single secret leaving the vault

## When it fails

- Exit code 1 on the first failure, with no retry
- Secrets written before the failure are left in place. Removing them is the caller's business

## License

This project is licensed under the [MIT License](./LICENSE).
