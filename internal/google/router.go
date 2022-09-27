package google

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

type Router struct {
	Name      string
	ProjectID string
	Region    string
	Network   string
}

func (r *Router) create(ctx context.Context, client *compute.RoutersClient) (*computepb.Router, error) {
	if r.exists(ctx, client) {
		return r.get(ctx, client)
	}
	req := &computepb.InsertRouterRequest{
		RouterResource: &computepb.Router{
			Name: &r.Name,
			Nats: []*computepb.RouterNat{
				{
					Name:                          &r.Name,
					NatIpAllocateOption:           StrPtr("AUTO_ONLY"),
					SourceSubnetworkIpRangesToNat: StrPtr("ALL_SUBNETWORKS_ALL_IP_RANGES"),
				},
			},
			Network: &r.Network,
		},
		Project: r.ProjectID,
		Region:  r.Region,
	}
	op, err := client.Insert(ctx, req)
	if err != nil {
		return nil, err
	}
	err = op.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return r.get(ctx, client)
}

func (r *Router) get(ctx context.Context, client *compute.RoutersClient) (*computepb.Router, error) {
	req := &computepb.GetRouterRequest{
		Project: r.ProjectID,
		Region:  r.Region,
		Router:  r.Name,
	}
	resp, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *Router) exists(ctx context.Context, client *compute.RoutersClient) bool {
	_, err := r.get(ctx, client)
	return err == nil
}

func (r *Router) delete(ctx context.Context, client *compute.RoutersClient) error {
	if r.exists(ctx, client) {
		req := &computepb.DeleteRouterRequest{
			Project: r.ProjectID,
			Region:  r.Region,
			Router:  r.Name,
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

func (r *Router) update(ctx context.Context, client *compute.RoutersClient) (*computepb.Router, error) {
	if r.exists(ctx, client) {
		return r.get(ctx, client)
	}
	req := &computepb.PatchRouterRequest{
		RouterResource: &computepb.Router{
			Name: &r.Name,
			Nats: []*computepb.RouterNat{
				{
					Name:                          &r.Name,
					NatIpAllocateOption:           StrPtr("AUTO_ONLY"),
					SourceSubnetworkIpRangesToNat: StrPtr("ALL_SUBNETWORKS_ALL_IP_RANGES"),
				},
			},
			Network: &r.Network,
		},
		Project: r.ProjectID,
		Region:  r.Region,
		Router:  r.Name,
	}
	op, err := client.Patch(ctx, req)
	if err != nil {
		return nil, err
	}
	err = op.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return r.get(ctx, client)
}
