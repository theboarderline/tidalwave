package google

import (
	"context"
	"fmt"
	"strings"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type CryptoKey struct {
	Name      string
	Keyring   string
	ProjectID string
}

func setIam(ctx context.Context, client *kms.KeyManagementClient, id string, key *kmspb.CryptoKey) error {
	proj, err := getProject(ctx, id)
	if err != nil {
		return err
	}
	member := fmt.Sprintf("serviceAccount:service-%s@container-engine-robot.iam.gserviceaccount.com", strings.Trim(proj.Name, "projects/"))
	handle := client.ResourceIAM(key.GetName())
	policy, err := handle.Policy(ctx)
	if err != nil {
		return err
	}
	policy.Add(member, "roles/cloudkms.cryptoKeyDecrypter")
	policy.Add(member, "roles/cloudkms.cryptoKeyEncrypter")
	return handle.SetPolicy(ctx, policy)
}

func removeIam(ctx context.Context, client *kms.KeyManagementClient, id string, key *kmspb.CryptoKey) error {
	proj, err := getProject(ctx, id)
	if err != nil {
		return err
	}
	member := fmt.Sprintf("serviceAccount:service-%s@container-engine-robot.iam.gserviceaccount.com", strings.Trim(proj.Name, "projects/"))
	handle := client.ResourceIAM(key.GetName())
	policy, err := handle.Policy(ctx)
	if err != nil {
		return err
	}
	policy.Remove(member, "roles/cloudkms.cryptoKeyDecrypter")
	policy.Remove(member, "roles/cloudkms.cryptoKeyEncrypter")
	return handle.SetPolicy(ctx, policy)
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

func (c *CryptoKey) delete(ctx context.Context, client *kms.KeyManagementClient) error {
	key, err := c.get(ctx, client)
	if err != nil {
		return err
	}
	return removeIam(ctx, client, c.ProjectID, key)
}
