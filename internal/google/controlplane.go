package google

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	container "cloud.google.com/go/container/apiv1"
	"github.com/kyokomi/emoji/v2"
)

type Controlplane struct {
	Vpc
	Subnetwork
	Router
	Cluster
}

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}

func (c *Controlplane) Create() error {
	ctx := context.Background()

	vpcClient, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer vpcClient.Close()
	network, err := c.Vpc.create(ctx, vpcClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane VPC created :check_mark_button:")

	c.Subnetwork.Network = network.GetSelfLink()
	subnetClient, err := compute.NewSubnetworksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer subnetClient.Close()
	_, err = c.Subnetwork.create(ctx, subnetClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane subnetwork created :check_mark_button:")

	c.Router.Network = network.GetSelfLink()
	routerClient, err := compute.NewRoutersRESTClient(ctx)
	if err != nil {
		return err
	}
	defer routerClient.Close()
	_, err = c.Router.create(ctx, routerClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane router created :check_mark_button:")

	clusterClient, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return err
	}
	defer clusterClient.Close()
	_, err = c.Cluster.create(ctx, clusterClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane cluster created :check_mark_button:")

	return nil
}

func (c *Controlplane) Delete() error {
	ctx := context.Background()

	clusterClient, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return err
	}
	defer clusterClient.Close()
	err = c.Cluster.delete(ctx, clusterClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane cluster destroyed :cross_mark_button:")

	routerClient, err := compute.NewRoutersRESTClient(ctx)
	if err != nil {
		return err
	}
	defer routerClient.Close()
	err = c.Router.delete(ctx, routerClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane router destroyed :cross_mark_button:")

	subnetClient, err := compute.NewSubnetworksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer subnetClient.Close()
	err = c.Subnetwork.delete(ctx, subnetClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane subnetwork destroyed :cross_mark_button:")

	vpcClient, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer vpcClient.Close()
	err = c.Vpc.delete(ctx, vpcClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane VPC destroyed :cross_mark_button:")

	return nil
}
