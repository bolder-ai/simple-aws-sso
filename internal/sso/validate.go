package sso

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type CallerIdentity struct {
	Account string
	Arn     string
	UserID  string
}

func ValidateCredentials(ctx context.Context, region string, creds *Credentials) (*CallerIdentity, error) {
	cfg := aws.Config{
		Region: region,
		Credentials: credentials.NewStaticCredentialsProvider(
			creds.AccessKeyID,
			creds.SecretAccessKey,
			creds.SessionToken,
		),
	}

	stsClient := sts.NewFromConfig(cfg)
	output, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("credentials validation failed: %w", err)
	}

	return &CallerIdentity{
		Account: aws.ToString(output.Account),
		Arn:     aws.ToString(output.Arn),
		UserID:  aws.ToString(output.UserId),
	}, nil
}
