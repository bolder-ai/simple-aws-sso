package awscreds

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"

	"simple-aws-sso/internal/sso"
)

const (
	awsAccessKeyID        = "aws_access_key_id"
	awsSecretAccessKey    = "aws_secret_access_key"
	awsSessionToken       = "aws_session_token"
	awsSecurityToken      = "aws_security_token"      // legacy alias for session token
	awsSessionExpiration  = "aws_session_expiration"  // ISO 8601 expiration time
)

func credentialsFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".aws", "credentials"), nil
}

func WriteCredentials(profile string, creds *sso.Credentials) error {
	credPath, err := credentialsFilePath()
	if err != nil {
		return err
	}

	// Ensure .aws directory exists
	awsDir := filepath.Dir(credPath)
	if err := os.MkdirAll(awsDir, 0700); err != nil {
		return fmt.Errorf("failed to create .aws directory: %w", err)
	}

	// Load existing credentials file or create new one
	cfg, err := ini.LooseLoad(credPath)
	if err != nil {
		return fmt.Errorf("failed to load credentials file: %w", err)
	}

	// Get or create the profile section
	section, err := cfg.NewSection(profile)
	if err != nil {
		// Section already exists, get it
		section = cfg.Section(profile)
	}

	// Set credentials
	section.Key(awsAccessKeyID).SetValue(creds.AccessKeyID)
	section.Key(awsSecretAccessKey).SetValue(creds.SecretAccessKey)
	section.Key(awsSessionToken).SetValue(creds.SessionToken)
	section.Key(awsSecurityToken).SetValue(creds.SessionToken)
	section.Key(awsSessionExpiration).SetValue(creds.Expiration.UTC().Format(time.RFC3339))

	// Write back to file
	if err := cfg.SaveTo(credPath); err != nil {
		return fmt.Errorf("failed to save credentials file: %w", err)
	}

	return nil
}
