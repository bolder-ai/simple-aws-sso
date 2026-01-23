package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"simple-aws-sso/internal/awscreds"
	"simple-aws-sso/internal/config"
	"simple-aws-sso/internal/sso"
	"simple-aws-sso/internal/ui"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-version" || os.Args[1] == "--version") {
		fmt.Println(version)
		os.Exit(0)
	}

	if err := run(); err != nil {
		ui.Error("%v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Authenticate via OIDC device flow
	oidcClient := sso.NewOIDCClient(cfg.Region)
	accessToken, err := oidcClient.Authenticate(ctx, cfg.SSOURL)
	if err != nil {
		return err
	}
	ui.Success("Authenticated")

	// Get role credentials
	credsClient := sso.NewCredentialsClient(cfg.Region)
	creds, err := credsClient.GetRoleCredentials(ctx, accessToken, cfg.AccountID, cfg.Role)
	if err != nil {
		return err
	}

	// Validate credentials if requested
	if cfg.ValidateCreds {
		ui.Wait("Validating credentials with STS...")
		identity, err := sso.ValidateCredentials(ctx, cfg.Region, creds)
		if err != nil {
			return err
		}
		ui.Success("Credentials validated")
		ui.Print("  %s %s\n", ui.Dim("ARN:"), identity.Arn)
	}

	// Write credentials to file
	if err := awscreds.WriteCredentials(cfg.Profile, creds); err != nil {
		return err
	}
	ui.Success("Credentials written to %s", ui.Cyan("~/.aws/credentials"))

	// Post-write validation: read back and verify
	if cfg.ValidateCreds {
		ui.Wait("Verifying saved credentials...")
		savedCreds, err := awscreds.ReadCredentials(cfg.Profile)
		if err != nil {
			return fmt.Errorf("failed to read back saved credentials: %w", err)
		}
		identity, err := sso.ValidateCredentials(ctx, cfg.Region, savedCreds)
		if err != nil {
			return fmt.Errorf("saved credentials are invalid: %w", err)
		}
		ui.Success("Saved credentials verified")
		ui.Print("  %s %s\n", ui.Dim("ARN:"), identity.Arn)
	}
	ui.Println()
	ui.Print("  %s  %s\n", ui.Dim("Profile:"), ui.Bold(cfg.Profile))
	ui.Print("  %s     %s\n", ui.Dim("Role:"), cfg.Role)
	ui.Print("  %s  %s\n", ui.Dim("Expires:"), ui.Yellow(creds.Expiration.Local().Format("2006-01-02 15:04:05")))

	return nil
}
