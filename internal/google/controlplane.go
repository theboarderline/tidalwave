package google

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
)

type Controlplane struct {
	Name      string
	ProjectID string
	Region    string
	Vpc
	Subnetwork
}

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}

func (c *Controlplane) Create() error {
	ctx := context.Background()
	c.Vpc = Vpc{
		Name:      c.Name,
		ProjectID: c.ProjectID,
	}
	vpcClient, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer vpcClient.Close()
	network, err := c.Vpc.create(ctx, vpcClient)
	if err != nil {
		return err
	}
	c.Subnetwork = Subnetwork{
		Name:      fmt.Sprintf("%s-controlplane", c.Name),
		ProjectID: c.ProjectID,
		Region:    c.Region,
		NodesCidr: c.Subnetwork.NodesCidr,
		PodsCidr:  c.Subnetwork.PodsCidr,
	}
	subnetClient, err := compute.NewSubnetworksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer subnetClient.Close()
	_, err = c.Subnetwork.create(ctx, subnetClient, network.SelfLink)
	if err != nil {
		return err
	}
	return nil
}
