# 🔐 simple-aws-sso

A minimal CLI tool to authenticate via AWS SSO and sync credentials to `~/.aws/credentials`.

## 📦 Installation

### Download binary

```bash
# Linux amd64
curl -sSL https://github.com/YOUR_USER/simple-aws-sso/releases/latest/download/simple-aws-sso_linux_amd64.tar.gz | tar xz

# Linux arm64
curl -sSL https://github.com/YOUR_USER/simple-aws-sso/releases/latest/download/simple-aws-sso_linux_arm64.tar.gz | tar xz

# macOS arm64 (M1/M2/M3)
curl -sSL https://github.com/YOUR_USER/simple-aws-sso/releases/latest/download/simple-aws-sso_darwin_arm64.tar.gz | tar xz

# macOS amd64
curl -sSL https://github.com/YOUR_USER/simple-aws-sso/releases/latest/download/simple-aws-sso_darwin_amd64.tar.gz | tar xz
```

### Build from source

```bash
go install github.com/YOUR_USER/simple-aws-sso@latest
```

## ⚙️ Configuration

Configuration is loaded with the following precedence (highest to lowest):

1. **Flags** — command line arguments
2. **Environment variables**
3. **~/.aws/config** — per-profile SSO settings

| Flag | Environment Variable | ~/.aws/config key |
|------|---------------------|-------------------|
| `--profile` | `AWS_PROFILE` | — |
| `--sso-url` | `AWS_SSO_URL` | `sso_start_url` |
| `--account-id` | `AWS_SSO_ACCOUNT_ID` | `sso_account_id` |
| `--region` | `AWS_REGION` | `sso_region` |
| `--role` | `AWS_SSO_ROLE` | `sso_role_name` |

## 🚀 Usage

### Using ~/.aws/config (simplest)

If your profile is already configured in `~/.aws/config`:

```ini
[profile dev]
sso_start_url = https://mycompany.awsapps.com/start
sso_region = eu-west-1
sso_account_id = 123456789012
sso_role_name = AdministratorAccess
```

Just run:

```bash
simple-aws-sso --profile dev
# or
AWS_PROFILE=dev simple-aws-sso
```

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

### Mixed (flags override env vars and config)

```bash
AWS_PROFILE=dev simple-aws-sso --role "ReadOnly"
```

## ✅ Output

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

## 🛠️ Development

Requires Go 1.25+ and [Task](https://taskfile.dev/).

```bash
task build             # Build to bin/
task run               # Build and run
task test              # Run tests
task clean             # Remove build artifacts
task release:snapshot  # Build release locally (requires goreleaser)
```

## 📄 License

MIT
