package main

import (
	"fmt"
	"syscall"

	"github.com/pkg/browser"
	"golang.org/x/term"
)

func (app *application) authorize() error {
	requestToken, _, err := app.config.RequestToken()
	if err != nil {
		return fmt.Errorf("failed to fetch request token: %w", err)
	}
	authorizationURL, err := app.config.AuthorizationURL(requestToken)
	if err != nil {
		return fmt.Errorf("failed to get authorization URL: %w", err)
	}
	err = browser.OpenURL(authorizationURL.String())
	if err != nil {
		return fmt.Errorf("failed to open URL '%s': %w", authorizationURL.String(), err)
	}

	fmt.Printf("Paste your PIN here: ")
	pin, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	accessToken, accessSecret, err := app.config.AccessToken(requestToken, "secret does not matter", string(pin))
	if err != nil {
		return fmt.Errorf("failed to fetch access token: %w", err)
	}

	fmt.Printf("\nSet the following environment variables or pass them with --token and --token-secret\n\n")
	fmt.Printf("SUPERBLOCK_TOKEN=%s\nSUPERBLOCK_TOKEN_SECRET=%s\n", accessToken, accessSecret)

	return nil
}
