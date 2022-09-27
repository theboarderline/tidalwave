package google

import (
	"context"
	"fmt"
	"strings"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
)

type CryptoKey struct {
	Name      string
	Keyring   string
	ProjectID string
	// Region    string
}

func setIam(ctx context.Context, client *kms.KeyManagementClient, id string, key *kmspb.CryptoKey) error {
	proj, err := getProject(ctx, id)
	if err != nil {
		return err
	}
	member := fmt.Sprintf("serviceAccount:service-%s@container-engine-robot.iam.gserviceaccount.com", strings.Trim(proj.Name, "projects/"))
	policy := &iampb.Policy{
		Bindings: []*iampb.Binding{
			{
				Role:    "roles/cloudkms.cryptoKeyDecrypter",
				Members: []string{member},
			},
			{
				Role:    "roles/cloudkms.cryptoKeyEncrypter",
				Members: []string{member},
			},
		},
	}
	iam := &iampb.SetIamPolicyRequest{
		Resource: key.Name,
		Policy:   policy,
	}
	_, err = client.SetIamPolicy(ctx, iam)
	if err != nil {
		return err
	}
	return nil
}

// Create KMS Crypto Key
func (c *CryptoKey) create(ctx context.Context, client *kms.KeyManagementClient) (*kmspb.CryptoKey, error) {
	if c.exists(ctx, client) {
		return c.get(ctx, client)
	}
	cryptoKey := &kmspb.CryptoKey{
		Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
	}
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:                     c.Keyring,
		CryptoKeyId:                c.Name,
		CryptoKey:                  cryptoKey,
		SkipInitialVersionCreation: false,
	}
	resp, err := client.CreateCryptoKey(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := setIam(ctx, client, c.ProjectID, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get KMS Crypto Key
func (c *CryptoKey) get(ctx context.Context, client *kms.KeyManagementClient) (*kmspb.CryptoKey, error) {
	req := &kmspb.GetCryptoKeyRequest{
		Name: fmt.Sprintf("%s/cryptoKeys/%s", c.Keyring, c.Name),
	}
	resp, err := client.GetCryptoKey(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Check if KMS Crypto Key exists
func (c *CryptoKey) exists(ctx context.Context, client *kms.KeyManagementClient) bool {
	_, err := c.get(ctx, client)
	return err == nil
}

// Update KMS Crypto Key
func (c *CryptoKey) update(ctx context.Context, client *kms.KeyManagementClient) (*kmspb.CryptoKey, error) {
	key, err := c.get(ctx, client)
	if err != nil {
		return nil, err
	}
	if err := setIam(ctx, client, c.ProjectID, key); err != nil {
		return nil, err
	}
	return key, nil
}
