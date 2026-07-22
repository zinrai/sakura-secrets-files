# sakura-secrets-pull

A command-line tool to pull secrets from Sakura Cloud Secret Manager and write them to local files.

## Features

- Pulls secrets from Sakura Cloud Secret Manager via API
- Writes secrets to specified file paths atomically
- Dry-run mode for testing

## Prerequisites

- Sakura Cloud account with Secret Manager enabled
- API credentials (Access Token and Access Token Secret)

## Configuration

### Environment Variables

```bash
$ export SAKURACLOUD_ACCESS_TOKEN="your-access-token"
$ export SAKURACLOUD_ACCESS_TOKEN_SECRET="your-access-token-secret"
```

### Configuration File

Create a YAML configuration file (e.g., `secrets-config.yaml`):

```yaml
vault:
  id: "123456789012"  # Your vault resource ID

secrets:
  - name: production-db-password
    dest: roles/database/files/db_password

  - name: production-api-key
    dest: roles/application/files/api_key
```

## Usage

### Basic Usage

```bash
$ sakura-secrets-pull -config secrets-config.yaml
```

### Dry-run Mode

Test what would be done without actually writing files:

```bash
$ sakura-secrets-pull -config secrets-config.yaml -dry-run
```

### Use Cases

This tool is useful for:

- Deployment automation scripts that need to fetch secrets before deployment
- CI/CD pipelines that manage secrets separately from code
- Configuration management workflows where secrets are stored centrally
- Local development environments that pull production-like secrets

Example deployment script:

```bash
#!/bin/bash
set -e

# Pull secrets from Secret Manager
sakura-secrets-pull -config secrets-production.yaml

# Use the pulled secrets in your deployment
# (e.g., configuration management tools, custom scripts, etc.)
./deploy.sh
```

## Command-line Options

- `-config <path>`: Path to configuration file (required)
- `-zone <zone>`: Sakura Cloud zone (default: is1a)
- `-dry-run`: Show what would be done without actually writing files

## Error Handling

The tool follows a fail-fast approach:

- If any secret fails to pull, the tool exits immediately with exit code 1
- No retry logic - re-run the command if it fails
- Errors include detailed messages about which operation failed

## Security

- Secrets are written with 0600 permissions (owner read/write only)
- Atomic writes using temporary files ensure no partial writes
- API credentials are never logged or written to disk

## License

This project is licensed under the [MIT License](./LICENSE).
