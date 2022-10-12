package google

import (
	"context"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type Keyring struct {
	Name      string
	ProjectID string
	Region    string
}

// Create KMS Keyring
func (k *Keyring) create(ctx context.Context, client *kms.KeyManagementClient) (*kmspb.KeyRing, error) {
	if k.exists(ctx, client) {
		return k.get(ctx, client)
	}
	keyring := &kmspb.KeyRing{}
	req := &kmspb.CreateKeyRingRequest{
		Parent:    fmt.Sprintf("projects/%s/locations/%s", k.ProjectID, k.Region),
		KeyRingId: k.Name,
		KeyRing:   keyring,
	}
	resp, err := client.CreateKeyRing(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get KMS Keyring
func (k *Keyring) get(ctx context.Context, client *kms.KeyManagementClient) (*kmspb.KeyRing, error) {
	req := &kmspb.GetKeyRingRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", k.ProjectID, k.Region, k.Name),
	}
	resp, err := client.GetKeyRing(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Check if KMS Keyring exists
func (k *Keyring) exists(ctx context.Context, client *kms.KeyManagementClient) bool {
	_, err := k.get(ctx, client)
	return err == nil
}
