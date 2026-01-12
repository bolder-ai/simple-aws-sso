# simple-aws-sso

A minimal CLI tool to authenticate via AWS SSO and sync credentials to `~/.aws/credentials`.

## Installation

### Download binary

```bash
# Linux amd64
curl -sSL https://github.com/YOUR_USER/simple-aws-sso/releases/download/v0.1.0/simple-aws-sso_0.1.0_linux_amd64.tar.gz | tar xz

# Linux arm64
curl -sSL https://github.com/YOUR_USER/simple-aws-sso/releases/download/v0.1.0/simple-aws-sso_0.1.0_linux_arm64.tar.gz | tar xz

# macOS arm64 (M1/M2/M3)
curl -sSL https://github.com/YOUR_USER/simple-aws-sso/releases/download/v0.1.0/simple-aws-sso_0.1.0_darwin_arm64.tar.gz | tar xz

# macOS amd64
curl -sSL https://github.com/YOUR_USER/simple-aws-sso/releases/download/v0.1.0/simple-aws-sso_0.1.0_darwin_amd64.tar.gz | tar xz
```

### Build from source

```bash
go install github.com/YOUR_USER/simple-aws-sso@latest
```

## Configuration

All options can be set via environment variables or flags. Flags take precedence.

| Flag | Environment Variable | Description |
|------|---------------------|-------------|
| `--sso-url` | `AWS_SSO_URL` | AWS SSO start URL |
| `--account-id` | `AWS_SSO_ACCOUNT_ID` | AWS account ID |
| `--profile` | `AWS_PROFILE` | Profile name in credentials file |
| `--region` | `AWS_REGION` | AWS region for SSO |
| `--role` | `AWS_SSO_ROLE` | AWS SSO role name |

## Usage

### With environment variables

```bash
export AWS_SSO_URL="https://mycompany.awsapps.com/start"
export AWS_SSO_ACCOUNT_ID="123456789012"
export AWS_SSO_ROLE="AdministratorAccess"
export AWS_PROFILE="dev"
export AWS_REGION="eu-west-1"

simple-aws-sso
```

### With flags

```bash
simple-aws-sso \
  --sso-url "https://mycompany.awsapps.com/start" \
  --account-id "123456789012" \
  --role "AdministratorAccess" \
  --profile "dev" \
  --region "eu-west-1"
```

### Mixed (flags override env vars)

```bash
export AWS_SSO_URL="https://mycompany.awsapps.com/start"
export AWS_REGION="eu-west-1"

simple-aws-sso --account-id "123456789012" --role "ReadOnly" --profile "prod"
```

## Output

```
→ Opening browser for SSO login...

  URL:  https://device.sso.eu-west-1.amazonaws.com/?user_code=ABCD-EFGH
  Code: ABCD-EFGH

◌ Waiting for authorization...
✓ Authenticated
✓ Credentials written to ~/.aws/credentials

  Profile:  dev
  Role:     AdministratorAccess
  Expires:  2026-01-12 14:30:00
```

## Development

Requires Go 1.25+ and [Task](https://taskfile.dev/).

```bash
task build          # Build to bin/
task run            # Build and run
task test           # Run tests
task clean          # Remove build artifacts
task release:snapshot  # Build release locally (requires goreleaser)
```

## License

MIT
