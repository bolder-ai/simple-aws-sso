package sso

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc/types"

	"simple-aws-sso/internal/ui"
)

const (
	clientName = "simple-aws-sso"
	clientType = "public"
)

type OIDCClient struct {
	client *ssooidc.Client
	region string
}

type DeviceAuth struct {
	UserCode        string
	VerificationURI string
	DeviceCode      string
	ExpiresIn       int32
	Interval        int32
}

func NewOIDCClient(region string) *OIDCClient {
	// SSO OIDC endpoints are unauthenticated - use anonymous credentials
	// to prevent the SDK from trying to sign requests
	cfg := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider("", "", ""),
	}
	return &OIDCClient{
		client: ssooidc.NewFromConfig(cfg),
		region: region,
	}
}

func (o *OIDCClient) Authenticate(ctx context.Context, startURL string) (string, error) {
	// Step 1: Register client
	regOutput, err := o.client.RegisterClient(ctx, &ssooidc.RegisterClientInput{
		ClientName: aws.String(clientName),
		ClientType: aws.String(clientType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to register client: %w", err)
	}

	// Step 2: Start device authorization
	authOutput, err := o.client.StartDeviceAuthorization(ctx, &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     regOutput.ClientId,
		ClientSecret: regOutput.ClientSecret,
		StartUrl:     aws.String(startURL),
	})
	if err != nil {
		return "", fmt.Errorf("failed to start device authorization: %w", err)
	}

	deviceAuth := &DeviceAuth{
		UserCode:        aws.ToString(authOutput.UserCode),
		VerificationURI: aws.ToString(authOutput.VerificationUriComplete),
		DeviceCode:      aws.ToString(authOutput.DeviceCode),
		ExpiresIn:       authOutput.ExpiresIn,
		Interval:        authOutput.Interval,
	}

	// Fallback to base URI if complete URI not provided
	if deviceAuth.VerificationURI == "" {
		deviceAuth.VerificationURI = aws.ToString(authOutput.VerificationUri)
	}

	// Open browser
	ui.Info("Opening browser for SSO login...")
	ui.Println()
	ui.Print("  %s  %s\n", ui.Dim("URL:"), ui.Cyan(deviceAuth.VerificationURI))
	ui.Print("  %s %s\n", ui.Dim("Code:"), ui.Bold(deviceAuth.UserCode))
	ui.Println()

	if err := openBrowser(deviceAuth.VerificationURI); err != nil {
		ui.Wait("Could not open browser automatically. Please open the URL manually.")
	}

	// Step 3: Poll for token
	ui.Wait("Waiting for authorization...")
	accessToken, err := o.pollForToken(ctx, regOutput, deviceAuth)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (o *OIDCClient) pollForToken(ctx context.Context, reg *ssooidc.RegisterClientOutput, auth *DeviceAuth) (string, error) {
	interval := time.Duration(auth.Interval) * time.Second
	if interval == 0 {
		interval = 5 * time.Second
	}

	timeout := time.Duration(auth.ExpiresIn) * time.Second
	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		tokenOutput, err := o.client.CreateToken(ctx, &ssooidc.CreateTokenInput{
			ClientId:     reg.ClientId,
			ClientSecret: reg.ClientSecret,
			DeviceCode:   aws.String(auth.DeviceCode),
			GrantType:    aws.String("urn:ietf:params:oauth:grant-type:device_code"),
		})

		if err != nil {
			var authPending *types.AuthorizationPendingException
			var slowDown *types.SlowDownException

			if errors.As(err, &authPending) {
				time.Sleep(interval)
				continue
			}
			if errors.As(err, &slowDown) {
				interval += 5 * time.Second
				time.Sleep(interval)
				continue
			}
			return "", fmt.Errorf("failed to create token: %w", err)
		}

		return aws.ToString(tokenOutput.AccessToken), nil
	}

	return "", errors.New("authorization timed out")
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
