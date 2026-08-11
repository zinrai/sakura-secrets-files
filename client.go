package main

import (
	"context"
	"fmt"
	"os"

	secretmanager "github.com/sacloud/sacloud-sdk-go/api/secretmanager"
	v1 "github.com/sacloud/sacloud-sdk-go/api/secretmanager/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

// SakuraClient represents the Sakura Cloud Secret Manager API client
type SakuraClient struct {
	Zone string
	api  *v1.Client
}

// NewClientFromEnv creates a new SakuraClient for the specified zone.
// Credentials are resolved by saclient-go from environment variables,
// supporting both static API keys and service principals.
func NewClientFromEnv(zone string) (*SakuraClient, error) {
	endpoint := fmt.Sprintf("SAKURA_ENDPOINTS_SECRETMANAGER=https://secure.sakura.ad.jp/cloud/zone/%s/api/cloud/1.1", zone)

	var sc saclient.Client
	if err := sc.SetEnviron(append(os.Environ(), endpoint)); err != nil {
		return nil, fmt.Errorf("failed to configure saclient: %w", err)
	}
	if err := sc.Populate(); err != nil {
		return nil, fmt.Errorf("failed to configure saclient: %w", err)
	}

	api, err := secretmanager.NewClient(&sc)
	if err != nil {
		return nil, fmt.Errorf("failed to create Secret Manager client: %w", err)
	}

	return &SakuraClient{
		Zone: zone,
		api:  api,
	}, nil
}

// GetSecret retrieves a secret from the specified vault using the unveil API.
// version 0 means the latest version.
func (c *SakuraClient) GetSecret(vaultID, secretName string, version int) (string, error) {
	op := secretmanager.NewSecretOp(c.api, vaultID)

	req := v1.Unveil{
		Name: secretName,
	}
	if version != 0 {
		req.Version = v1.NewOptNilInt(version)
	}

	res, err := op.Unveil(context.Background(), req)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}
