package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

type Config struct {
	SSOURL    string
	AccountID string
	Profile   string
	Region    string
	Role      string
}

// awsConfigProfile holds SSO settings from ~/.aws/config
type awsConfigProfile struct {
	SSOStartURL string
	SSORegion   string
	SSOAccountID string
	SSORoleName  string
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Parse flags
	flag.StringVar(&cfg.SSOURL, "sso-url", "", "AWS SSO start URL")
	flag.StringVar(&cfg.AccountID, "account-id", "", "AWS account ID")
	flag.StringVar(&cfg.Profile, "profile", "", "AWS profile name for credentials file")
	flag.StringVar(&cfg.Region, "region", "", "AWS region for SSO")
	flag.StringVar(&cfg.Role, "role", "", "AWS SSO role name")
	flag.Parse()

	// Get profile first (needed to load aws config)
	// Precedence: flag > env var
	if cfg.Profile == "" {
		cfg.Profile = os.Getenv("AWS_PROFILE")
	}

	// Load AWS config profile settings as lowest precedence
	var awsCfg *awsConfigProfile
	if cfg.Profile != "" {
		awsCfg = loadAWSConfigProfile(cfg.Profile)
	}

	// Apply precedence: flag > env var > aws config
	if cfg.SSOURL == "" {
		cfg.SSOURL = os.Getenv("AWS_SSO_URL")
	}
	if cfg.SSOURL == "" && awsCfg != nil {
		cfg.SSOURL = awsCfg.SSOStartURL
	}

	if cfg.AccountID == "" {
		cfg.AccountID = os.Getenv("AWS_SSO_ACCOUNT_ID")
	}
	if cfg.AccountID == "" && awsCfg != nil {
		cfg.AccountID = awsCfg.SSOAccountID
	}

	if cfg.Region == "" {
		cfg.Region = os.Getenv("AWS_REGION")
	}
	if cfg.Region == "" && awsCfg != nil {
		cfg.Region = awsCfg.SSORegion
	}

	if cfg.Role == "" {
		cfg.Role = os.Getenv("AWS_SSO_ROLE")
	}
	if cfg.Role == "" && awsCfg != nil {
		cfg.Role = awsCfg.SSORoleName
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadAWSConfigProfile(profile string) *awsConfigProfile {
	configPath, err := awsConfigFilePath()
	if err != nil {
		return nil
	}

	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil
	}

	// AWS config uses "profile <name>" for named profiles, but "default" for default
	sectionName := fmt.Sprintf("profile %s", profile)
	if profile == "default" {
		sectionName = "default"
	}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return nil
	}

	return &awsConfigProfile{
		SSOStartURL:  section.Key("sso_start_url").String(),
		SSORegion:    section.Key("sso_region").String(),
		SSOAccountID: section.Key("sso_account_id").String(),
		SSORoleName:  section.Key("sso_role_name").String(),
	}
}

func awsConfigFilePath() (string, error) {
	if configFile := os.Getenv("AWS_CONFIG_FILE"); configFile != "" {
		return configFile, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".aws", "config"), nil
}

func (c *Config) Validate() error {
	var missing []string

	if c.SSOURL == "" {
		missing = append(missing, "sso_start_url (--sso-url, AWS_SSO_URL, or ~/.aws/config)")
	}
	if c.AccountID == "" {
		missing = append(missing, "sso_account_id (--account-id, AWS_SSO_ACCOUNT_ID, or ~/.aws/config)")
	}
	if c.Profile == "" {
		missing = append(missing, "profile (--profile or AWS_PROFILE)")
	}
	if c.Region == "" {
		missing = append(missing, "sso_region (--region, AWS_REGION, or ~/.aws/config)")
	}
	if c.Role == "" {
		missing = append(missing, "sso_role_name (--role, AWS_SSO_ROLE, or ~/.aws/config)")
	}

	if len(missing) > 0 {
		return errors.New("missing required configuration:\n  - " + strings.Join(missing, "\n  - "))
	}

	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("SSO URL: %s, Account: %s, Profile: %s, Region: %s, Role: %s",
		c.SSOURL, c.AccountID, c.Profile, c.Region, c.Role)
}
