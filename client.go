package main

import (
	"context"
	"fmt"
	"os"

	secretmanager "github.com/sacloud/sacloud-sdk-go/api/secretmanager"
	v1 "github.com/sacloud/sacloud-sdk-go/api/secretmanager/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

type SakuraClient struct {
	api *v1.Client
}

// Credentials are left to saclient rather than read here, so that API keys and
// service principals keep resolving the way the SDK documents them
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

	return &SakuraClient{api: api}, nil
}

func (c *SakuraClient) GetSecret(vaultID, secretName string) (string, error) {
	op := secretmanager.NewSecretOp(c.api, vaultID)

	req := v1.Unveil{
		Name: secretName,
	}

	res, err := op.Unveil(context.Background(), req)
	if err != nil {
		return "", err
	}

	return res.Value, nil
}
