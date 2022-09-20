package google

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

type Vpc struct {
	Name      string
	ProjectID string
}

func (n *Vpc) create(ctx context.Context, client *compute.NetworksClient) (*computepb.Network, error) {
	if n.exists(ctx, client) {
		return n.get(ctx, client)
	}
	req := &computepb.InsertNetworkRequest{
		NetworkResource: &computepb.Network{
			AutoCreateSubnetworks: boolPtr(false),
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

func (n *Vpc) exists(ctx context.Context, client *compute.NetworksClient) bool {
	_, err := n.get(ctx, client)
	return err == nil
}

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
