package google

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	container "cloud.google.com/go/container/apiv1"
	kms "cloud.google.com/go/kms/apiv1"
	"github.com/kyokomi/emoji/v2"
)

type Controlplane struct {
	Vpc
	Subnetwork
	Router
	Cluster
	Firewalls []Firewall
	Keyring
	CryptoKey
}

func BoolPtr(b bool) *bool {
	return &b
}

func StrPtr(s string) *string {
	return &s
}

// Create controlplane
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

	kmsClient, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	defer kmsClient.Close()

	keyring, err := c.Keyring.create(ctx, kmsClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane KMS Keyring created :check_mark_button:")

	c.CryptoKey.Keyring = keyring.Name
	cryptoKey, err := c.CryptoKey.create(ctx, kmsClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane KMS Crypto Key created :check_mark_button:")

	c.Cluster.CryptoKeyName = cryptoKey.Name
	// clusterClient, err := container.NewClusterManagerClient(ctx)
	// if err != nil {
	// 	return err
	// }
	// defer clusterClient.Close()
	// _, err = c.Cluster.create(ctx, clusterClient)
	// if err != nil {
	// 	return err
	// }
	// emoji.Println("Controlplane cluster created :check_mark_button:")

	// firewallClient, err := compute.NewFirewallsRESTClient(ctx)
	// if err != nil {
	// 	return err
	// }
	// defer firewallClient.Close()
	// for i := range c.Firewalls {
	// 	c.Firewalls[i].Network = network.GetSelfLink()
	// 	_, err = c.Firewalls[i].create(ctx, firewallClient)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// emoji.Println("Controlplane firewall rules created :check_mark_button:")

	return nil
}

// Delete controlplane
func (c *Controlplane) Delete() error {
	ctx := context.Background()

	firewallClient, err := compute.NewFirewallsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer firewallClient.Close()
	for i := range c.Firewalls {
		err = c.Firewalls[i].delete(ctx, firewallClient)
		if err != nil {
			return err
		}
	}
	emoji.Println("Controlplane firewall rules deleted :check_mark_button:")

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

// Update controlplane
func (c *Controlplane) Update() error {
	ctx := context.Background()

	vpcClient, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer vpcClient.Close()
	network, err := c.Vpc.update(ctx, vpcClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane VPC updated :check_mark_button:")

	c.Subnetwork.Network = network.GetSelfLink()
	subnetClient, err := compute.NewSubnetworksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer subnetClient.Close()
	_, err = c.Subnetwork.update(ctx, subnetClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane subnetwork updated :check_mark_button:")

	c.Router.Network = network.GetSelfLink()
	routerClient, err := compute.NewRoutersRESTClient(ctx)
	if err != nil {
		return err
	}
	defer routerClient.Close()
	_, err = c.Router.update(ctx, routerClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane router updated :check_mark_button:")

	kmsClient, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	defer kmsClient.Close()
	keyring, err := c.Keyring.get(ctx, kmsClient)
	if err != nil {
		return err
	}
	c.CryptoKey.Keyring = keyring.Name
	_, err = c.CryptoKey.update(ctx, kmsClient)
	if err != nil {
		return err
	}
	emoji.Println("Controlplane KMS Crypto Key updated :check_mark_button:")

	// clusterClient, err := container.NewClusterManagerClient(ctx)
	// if err != nil {
	// 	return err
	// }
	// defer clusterClient.Close()
	// _, err = c.Cluster.update(ctx, clusterClient)
	// if err != nil {
	// 	return err
	// }
	// emoji.Println("Controlplane cluster updated :check_mark_button:")

	// firewallClient, err := compute.NewFirewallsRESTClient(ctx)
	// if err != nil {
	// 	return err
	// }
	// defer firewallClient.Close()
	// for i := range c.Firewalls {
	// 	c.Firewalls[i].Network = network.GetSelfLink()
	// 	_, err = c.Firewalls[i].update(ctx, firewallClient)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// emoji.Println("Controlplane firewall rules updated :check_mark_button:")

	return nil
}
