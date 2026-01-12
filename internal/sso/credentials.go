package sso

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sso"
)

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

type CredentialsClient struct {
	client *sso.Client
}

func NewCredentialsClient(region string) *CredentialsClient {
	// SSO GetRoleCredentials uses bearer token auth, not AWS credentials
	cfg := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider("", "", ""),
	}
	return &CredentialsClient{
		client: sso.NewFromConfig(cfg),
	}
}

func (c *CredentialsClient) GetRoleCredentials(ctx context.Context, accessToken, accountID, roleName string) (*Credentials, error) {
	output, err := c.client.GetRoleCredentials(ctx, &sso.GetRoleCredentialsInput{
		AccessToken: aws.String(accessToken),
		AccountId:   aws.String(accountID),
		RoleName:    aws.String(roleName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get role credentials: %w", err)
	}

	creds := output.RoleCredentials
	return &Credentials{
		AccessKeyID:     aws.ToString(creds.AccessKeyId),
		SecretAccessKey: aws.ToString(creds.SecretAccessKey),
		SessionToken:    aws.ToString(creds.SessionToken),
		Expiration:      time.UnixMilli(creds.Expiration),
	}, nil
}
