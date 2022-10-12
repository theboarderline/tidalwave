package google

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

type Firewall struct {
	Name              string
	ProjectID         string
	Allowed           []*computepb.Allowed
	DestinationRanges []string
	Direction         string
	Network           string
	SourceRanges      []string
	SourceTags        []string
	TargetTags        []string
}

// Create firewall rule
func (f *Firewall) create(ctx context.Context, client *compute.FirewallsClient) (*computepb.Firewall, error) {
	if f.exists(ctx, client) {
		return f.get(ctx, client)
	}

	req := &computepb.InsertFirewallRequest{
		FirewallResource: &computepb.Firewall{
			Allowed:           f.Allowed,
			DestinationRanges: f.DestinationRanges,
			Direction:         StrPtr(f.Direction),
			Name:              StrPtr(f.Name),
			Network:           StrPtr(f.Network),
			SourceRanges:      f.SourceRanges,
			SourceTags:        f.SourceTags,
			TargetTags:        f.TargetTags,
		},
		Project: f.ProjectID,
	}

	op, err := client.Insert(ctx, req)

	if err != nil {
		return nil, err
	}

	err = op.Wait(ctx)

	if err != nil {
		return nil, err
	}

	return f.get(ctx, client)
}

// Get firewall rule
func (f *Firewall) get(ctx context.Context, client *compute.FirewallsClient) (*computepb.Firewall, error) {
	req := &computepb.GetFirewallRequest{
		Firewall: f.Name,
		Project:  f.ProjectID,
	}
	resp, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Check if firewall rule exists
func (f *Firewall) exists(ctx context.Context, client *compute.FirewallsClient) bool {
	_, err := f.get(ctx, client)
	return err == nil
}

// Delete firewall rule
func (f *Firewall) delete(ctx context.Context, client *compute.FirewallsClient) error {
	if f.exists(ctx, client) {
		req := &computepb.DeleteFirewallRequest{
			Firewall: f.Name,
			Project:  f.ProjectID,
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

// Update firewall rule
func (f *Firewall) update(ctx context.Context, client *compute.FirewallsClient) (*computepb.Firewall, error) {
	if f.exists(ctx, client) {
		return f.get(ctx, client)
	}

	req := &computepb.PatchFirewallRequest{
		FirewallResource: &computepb.Firewall{
			Allowed:           f.Allowed,
			DestinationRanges: f.DestinationRanges,
			Direction:         StrPtr(f.Direction),
			Name:              StrPtr(f.Name),
			Network:           StrPtr(f.Network),
			SourceRanges:      f.SourceRanges,
			SourceTags:        f.SourceTags,
			TargetTags:        f.TargetTags,
		},
		Project:  f.ProjectID,
		Firewall: f.Name,
	}

	op, err := client.Patch(ctx, req)

	if err != nil {
		return nil, err
	}

	err = op.Wait(ctx)

	if err != nil {
		return nil, err
	}

	return f.get(ctx, client)
}
