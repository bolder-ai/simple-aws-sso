package awscreds

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"simple-aws-sso/internal/sso"
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

	// Read existing file content
	existingContent := ""
	if data, err := os.ReadFile(credPath); err == nil {
		existingContent = string(data)
	}

	// Remove existing profile section if present
	existingContent = removeProfileSection(existingContent, profile)

	// Build new profile section (matching yawsso format exactly)
	// yawsso uses ConfigParser which writes "key = value" with spaces
	newSection := fmt.Sprintf("[%s]\n", profile)
	newSection += fmt.Sprintf("aws_access_key_id = %s\n", creds.AccessKeyID)
	newSection += fmt.Sprintf("aws_secret_access_key = %s\n", creds.SecretAccessKey)
	newSection += fmt.Sprintf("aws_session_token = %s\n", creds.SessionToken)
	newSection += fmt.Sprintf("aws_security_token = %s\n", creds.SessionToken) // legacy alias
	newSection += fmt.Sprintf("aws_session_expiration = %s\n", creds.Expiration.UTC().Format("2006-01-02T15:04:05+0000"))

	// Append new section
	finalContent := strings.TrimRight(existingContent, "\n")
	if finalContent != "" {
		finalContent += "\n\n"
	}
	finalContent += newSection

	// Write back to file
	if err := os.WriteFile(credPath, []byte(finalContent), 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

// removeProfileSection removes a profile section from the credentials content
func removeProfileSection(content, profile string) string {
	lines := strings.Split(content, "\n")
	var result []string
	sectionHeader := fmt.Sprintf("[%s]", profile)
	inTargetSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if we're entering the target section
		if trimmed == sectionHeader {
			inTargetSection = true
			continue
		}

		// Check if we're entering a new section (leaving target section)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			inTargetSection = false
		}

		// Keep lines that aren't in the target section
		if !inTargetSection {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// ReadCredentials reads credentials for a specific profile (for validation)
func ReadCredentials(profile string) (*sso.Credentials, error) {
	credPath, err := credentialsFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(credPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open credentials file: %w", err)
	}
	defer file.Close()

	creds := &sso.Credentials{}
	inSection := false
	sectionHeader := fmt.Sprintf("[%s]", profile)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == sectionHeader {
			inSection = true
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if inSection {
				break // We've left our section
			}
			continue
		}

		if inSection && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "aws_access_key_id":
				creds.AccessKeyID = value
			case "aws_secret_access_key":
				creds.SecretAccessKey = value
			case "aws_session_token":
				creds.SessionToken = value
			case "aws_session_expiration":
				if t, err := time.Parse(time.RFC3339, value); err == nil {
					creds.Expiration = t
				}
			}
		}
	}

	if creds.AccessKeyID == "" {
		return nil, fmt.Errorf("profile [%s] not found", profile)
	}

	return creds, nil
}
