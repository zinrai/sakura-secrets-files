package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// SakuraClient represents the Sakura Cloud API client
type SakuraClient struct {
	Token       string
	TokenSecret string
	Zone        string
	HTTPClient  *http.Client
}

// UnveilRequest represents the request body for the unveil API
type UnveilRequest struct {
	Secret SecretRequest `json:"Secret"`
}

// SecretRequest represents the secret request parameters
type SecretRequest struct {
	Name    string `json:"Name"`
	Version int    `json:"Version"`
}

// UnveilResponse represents the response from the unveil API
type UnveilResponse struct {
	Secret SecretResponse `json:"Secret"`
}

// SecretResponse represents the secret response data
type SecretResponse struct {
	Name    string `json:"Name"`
	Version int    `json:"Version"`
	Value   string `json:"Value"`
}

// NewClientFromEnv creates a new SakuraClient from environment variables
func NewClientFromEnv(zone string) (*SakuraClient, error) {
	token := os.Getenv("SAKURACLOUD_ACCESS_TOKEN")
	tokenSecret := os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET")

	if token == "" {
		return nil, fmt.Errorf("SAKURACLOUD_ACCESS_TOKEN is not set")
	}
	if tokenSecret == "" {
		return nil, fmt.Errorf("SAKURACLOUD_ACCESS_TOKEN_SECRET is not set")
	}

	return &SakuraClient{
		Token:       token,
		TokenSecret: tokenSecret,
		Zone:        zone,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// GetSecret retrieves a secret from the specified vault
func (c *SakuraClient) GetSecret(vaultID, secretName string, version int) (string, error) {
	url := fmt.Sprintf(
		"https://secure.sakura.ad.jp/cloud/zone/%s/api/cloud/1.1/secretmanager/vaults/%s/secrets/unveil",
		c.Zone,
		vaultID,
	)

	reqBody := UnveilRequest{
		Secret: SecretRequest{
			Name:    secretName,
			Version: version,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Token, c.TokenSecret)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var unveilResp UnveilResponse
	if err := json.Unmarshal(body, &unveilResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return unveilResp.Secret.Value, nil
}
