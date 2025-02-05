package google

import (
	"context"
	"fmt"
	"log"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// CryptoKey represents a KMS Crypto Key
type CryptoKey struct {
	Name          string
	Keyring       string
	ProjectID     string
	ProjectNumber string
}

func setIam(ctx context.Context, client *kms.KeyManagementClient, id string, key *kmspb.CryptoKey) error {
	member := fmt.Sprintf("serviceAccount:service-%s@container-engine-robot.iam.gserviceaccount.com", id)
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
	member := fmt.Sprintf("serviceAccount:service-%s@container-engine-robot.iam.gserviceaccount.com", id)
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
	k, ok := c.exists(ctx, client)
	if ok {
		if err := setIam(ctx, client, c.ProjectNumber, k); err != nil {
			return nil, err
		}
		_, err := c.checkVersion(ctx, client, k)
		if err != nil {
			return nil, err
		}
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
	if err := setIam(ctx, client, c.ProjectNumber, resp); err != nil {
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
func (c *CryptoKey) exists(ctx context.Context, client *kms.KeyManagementClient) (*kmspb.CryptoKey, bool) {
	k, err := c.get(ctx, client)
	if err != nil {
		return nil, false
	}
	return k, err == nil
}

func (c *CryptoKey) checkVersion(ctx context.Context, client *kms.KeyManagementClient, k *kmspb.CryptoKey) (*kmspb.CryptoKeyVersion, error) {
	kv := k.GetPrimary()
	switch kv.State {
	case kmspb.CryptoKeyVersion_DISABLED:
		log.Printf("cryptokey version %s is disabled\n", k.GetName())
		ks, err := enableKeyVersion(ctx, client, kv)
		if err != nil {
			return ks, err
		}
		return ks, nil
	case kmspb.CryptoKeyVersion_DESTROYED:
		log.Printf("cryptokey %s is destroyed\n", k.GetName())
		ks, err := createKeyVersion(ctx, client, k)
		if err != nil {
			return ks, err
		}
	case kmspb.CryptoKeyVersion_DESTROY_SCHEDULED:
		log.Printf("cryptokey %s is scheduled to be destroyed\n", k.GetName())
		ks, err := restoreKeyVersion(ctx, client, kv)
		if err != nil {
			return ks, err
		}
		ks, err = enableKeyVersion(ctx, client, ks)
		if err != nil {
			return ks, err
		}
		return ks, nil
	}
	return kv, nil
}

// Update KMS Crypto Key
func (c *CryptoKey) update(ctx context.Context, client *kms.KeyManagementClient) (*kmspb.CryptoKey, error) {
	key, err := c.get(ctx, client)
	if err != nil {
		return nil, err
	}
	if err := setIam(ctx, client, c.ProjectNumber, key); err != nil {
		return nil, err
	}
	return key, nil
}

func (c *CryptoKey) delete(ctx context.Context, client *kms.KeyManagementClient) error {
	key, err := c.get(ctx, client)
	if err != nil {
		return err
	}
	return removeIam(ctx, client, c.ProjectNumber, key)
}

func enableKeyVersion(ctx context.Context, client *kms.KeyManagementClient, k *kmspb.CryptoKeyVersion) (*kmspb.CryptoKeyVersion, error) {
	req := &kmspb.UpdateCryptoKeyVersionRequest{
		CryptoKeyVersion: &kmspb.CryptoKeyVersion{
			Name:  k.GetName(),
			State: kmspb.CryptoKeyVersion_ENABLED,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"state"},
		},
	}
	kv, err := client.UpdateCryptoKeyVersion(ctx, req)
	if err != nil {
		return nil, err
	}
	var ks *kmspb.CryptoKeyVersion
status:
	for {
		ks, err := client.GetCryptoKeyVersion(ctx, &kmspb.GetCryptoKeyVersionRequest{
			Name: kv.GetName(),
		})
		if err != nil {
			return nil, err
		}
		if ks.GetState() == kmspb.CryptoKeyVersion_ENABLED {
			break status
		}
	}
	return ks, nil
}

func createKeyVersion(ctx context.Context, client *kms.KeyManagementClient, k *kmspb.CryptoKey) (*kmspb.CryptoKeyVersion, error) {
	req := &kmspb.CreateCryptoKeyVersionRequest{
		Parent: k.GetName(),
	}
	kv, err := client.CreateCryptoKeyVersion(ctx, req)
	if err != nil {
		return nil, err
	}
	updateReq := &kmspb.UpdateCryptoKeyPrimaryVersionRequest{
		Name:               k.GetName(),
		CryptoKeyVersionId: kv.GetName(),
	}
	_, err = client.UpdateCryptoKeyPrimaryVersion(ctx, updateReq)
	if err != nil {
		return nil, err
	}
	return kv, nil
}

func restoreKeyVersion(ctx context.Context, client *kms.KeyManagementClient, k *kmspb.CryptoKeyVersion) (*kmspb.CryptoKeyVersion, error) {
	req := &kmspb.RestoreCryptoKeyVersionRequest{
		Name: k.GetName(),
	}
	kv, err := client.RestoreCryptoKeyVersion(ctx, req)
	if err != nil {
		return nil, err
	}
status:
	for {
		ks, err := client.GetCryptoKeyVersion(ctx, &kmspb.GetCryptoKeyVersionRequest{
			Name: kv.GetName(),
		})
		if err != nil {
			return nil, err
		}
		if ks.GetState() == kmspb.CryptoKeyVersion_DISABLED {
			break status
		}
	}
	return kv, nil
}
