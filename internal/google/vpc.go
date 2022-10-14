package google

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// Vpc represents a VPC
type Vpc struct {
	Name      string
	ProjectID string
}

// Create VPC
func (n *Vpc) create(ctx context.Context, client *compute.NetworksClient) (*computepb.Network, error) {
	if n.exists(ctx, client) {
		return n.get(ctx, client)
	}
	req := &computepb.InsertNetworkRequest{
		NetworkResource: &computepb.Network{
			AutoCreateSubnetworks: BoolPtr(false),
			Name:                  &n.Name,
		},
		Project: n.ProjectID,
	}
	op, err := client.Insert(ctx, req)
	if err != nil {
		return nil, err
	}
	err = op.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return n.get(ctx, client)
}

// Get VPC
func (n *Vpc) get(ctx context.Context, client *compute.NetworksClient) (*computepb.Network, error) {
	req := &computepb.GetNetworkRequest{
		Project: n.ProjectID,
		Network: n.Name,
	}
	resp, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Check if VPC exists
func (n *Vpc) exists(ctx context.Context, client *compute.NetworksClient) bool {
	_, err := n.get(ctx, client)
	return err == nil
}

// Delete VPC
func (n *Vpc) delete(ctx context.Context, client *compute.NetworksClient) error {
	if n.exists(ctx, client) {
		req := &computepb.DeleteNetworkRequest{
			Project: n.ProjectID,
			Network: n.Name,
		}
		op, err := client.Delete(ctx, req)
		if err != nil {
			return err
		}
		err = op.Wait(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update VPC
func (n *Vpc) update(ctx context.Context, client *compute.NetworksClient) (*computepb.Network, error) {
	if n.exists(ctx, client) {
		return n.get(ctx, client)
	}

	req := &computepb.PatchNetworkRequest{
		Network: n.Name,
		NetworkResource: &computepb.Network{
			AutoCreateSubnetworks: BoolPtr(false),
			Name:                  &n.Name,
		},
		Project: n.ProjectID,
	}
	op, err := client.Patch(ctx, req)
	if err != nil {
		return nil, err
	}
	err = op.Wait(ctx)
	if err != nil {
		return nil, err
	}

	return n.get(ctx, client)
}
