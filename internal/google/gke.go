package google

import (
	"context"
	"errors"
	"fmt"

	container "cloud.google.com/go/container/apiv1"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

type Cluster struct {
	Name                 string
	ProjectID            string
	Region               string
	Network              string
	Subnetwork           string
	MachineType          string
	DiskSizeGb           int32
	MinNodeCount         int32
	MaxNodeCount         int32
	MasterAuthCidrBlocks []*containerpb.MasterAuthorizedNetworksConfig_CidrBlock
	MasterIpv4CidrBlock  string
}

func (c *Cluster) create(ctx context.Context, client *container.ClusterManagerClient) (*containerpb.Cluster, error) {
	if c.exists(ctx, client) {
		return c.get(ctx, client)
	}

	req := &containerpb.CreateClusterRequest{
		Cluster: &containerpb.Cluster{
			Name:    c.Name,
			Network: c.Network,
			AddonsConfig: &containerpb.AddonsConfig{
				HttpLoadBalancing: &containerpb.HttpLoadBalancing{
					Disabled: false,
				},
				HorizontalPodAutoscaling: &containerpb.HorizontalPodAutoscaling{
					Disabled: false,
				},
				ConfigConnectorConfig: &containerpb.ConfigConnectorConfig{
					Enabled: true,
				},
				GcePersistentDiskCsiDriverConfig: &containerpb.GcePersistentDiskCsiDriverConfig{
					Enabled: true,
				},
				GcpFilestoreCsiDriverConfig: &containerpb.GcpFilestoreCsiDriverConfig{
					Enabled: true,
				},
			},
			Subnetwork: c.Subnetwork,
			NodePools: []*containerpb.NodePool{
				{
					Name: "default-pool",
					Config: &containerpb.NodeConfig{
						MachineType: c.MachineType,
						DiskSizeGb:  c.DiskSizeGb,
						OauthScopes: []string{
							"https://www.googleapis.com/auth/devstorage.read_only",
							"https://www.googleapis.com/auth/logging.write",
							"https://www.googleapis.com/auth/monitoring",
							"https://www.googleapis.com/auth/servicecontrol",
							"https://www.googleapis.com/auth/service.management.readonly",
							"https://www.googleapis.com/auth/trace.append",
							"https://www.googleapis.com/auth/cloud-platform",
						},
						Tags: []string{
							fmt.Sprintf("%s-controlplane-default-pool", c.Name),
						},
						DiskType: "pd-ssd",
						WorkloadMetadataConfig: &containerpb.WorkloadMetadataConfig{
							Mode: 2,
						},
						ShieldedInstanceConfig: &containerpb.ShieldedInstanceConfig{
							EnableSecureBoot: true,
						},
					},
					InitialNodeCount: 1,
					Autoscaling: &containerpb.NodePoolAutoscaling{
						Enabled:      true,
						MinNodeCount: c.MinNodeCount,
						MaxNodeCount: c.MaxNodeCount,
					},
					Management: &containerpb.NodeManagement{
						AutoUpgrade: true,
						AutoRepair:  true,
					},
					UpgradeSettings: &containerpb.NodePool_UpgradeSettings{
						MaxSurge:       1,
						MaxUnavailable: 1,
					},
				},
			},
			IpAllocationPolicy: &containerpb.IPAllocationPolicy{
				UseIpAliases:               true,
				ClusterSecondaryRangeName:  "pods",
				ServicesSecondaryRangeName: "services",
			},
			MasterAuthorizedNetworksConfig: &containerpb.MasterAuthorizedNetworksConfig{
				Enabled:    true,
				CidrBlocks: c.MasterAuthCidrBlocks,
			},
			MaintenancePolicy: &containerpb.MaintenancePolicy{
				Window: &containerpb.MaintenanceWindow{
					Policy: &containerpb.MaintenanceWindow_DailyMaintenanceWindow{
						DailyMaintenanceWindow: &containerpb.DailyMaintenanceWindow{
							StartTime: "06:00",
						},
					},
				},
			},
			BinaryAuthorization: &containerpb.BinaryAuthorization{
				Enabled: true,
			},
			NetworkConfig: &containerpb.NetworkConfig{
				EnableIntraNodeVisibility: true,
				DatapathProvider:          0,
			},
			PrivateClusterConfig: &containerpb.PrivateClusterConfig{
				EnablePrivateNodes:  true,
				MasterIpv4CidrBlock: c.MasterIpv4CidrBlock,
			},
			ShieldedNodes: &containerpb.ShieldedNodes{
				Enabled: true,
			},
			ReleaseChannel: &containerpb.ReleaseChannel{
				Channel: 1,
			},
			WorkloadIdentityConfig: &containerpb.WorkloadIdentityConfig{
				WorkloadPool: fmt.Sprintf("%s.svc.id.goog", c.ProjectID),
			},
		},
		Parent: fmt.Sprintf("projects/%s/locations/%s", c.ProjectID, c.Region),
	}

	op, err := client.CreateCluster(ctx, req)

	if err != nil {
		return nil, err
	}

status:
	for {
		s, err := client.GetOperation(ctx, &containerpb.GetOperationRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/operations/%s", c.ProjectID, c.Region, op.GetName()),
		})
		if err != nil {
			return nil, err
		}
		switch s.GetStatus().Number() {
		case 3:
			break status
		case 4:
			return nil, errors.New(s.GetError().Message)
		}
	}
	return c.get(ctx, client)
}

func (c *Cluster) get(ctx context.Context, client *container.ClusterManagerClient) (*containerpb.Cluster, error) {
	req := &containerpb.GetClusterRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", c.ProjectID, c.Region, c.Name),
	}

	resp, err := client.GetCluster(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Cluster) exists(ctx context.Context, client *container.ClusterManagerClient) bool {
	_, err := c.get(ctx, client)
	return err == nil
}

func (c *Cluster) delete(ctx context.Context, client *container.ClusterManagerClient) error {
	if c.exists(ctx, client) {
		req := &containerpb.DeleteClusterRequest{
			Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", c.ProjectID, c.Region, c.Name),
		}
		op, err := client.DeleteCluster(ctx, req)
		if err != nil {
			return err
		}
	status:
		for {
			s, err := client.GetOperation(ctx, &containerpb.GetOperationRequest{
				Name: fmt.Sprintf("projects/%s/locations/%s/operations/%s", c.ProjectID, c.Region, op.GetName()),
			})
			if err != nil {
				return err
			}
			switch s.GetStatus().Number() {
			case 3:
				break status
			case 4:
				return errors.New(s.GetError().Message)
			}
		}
	}
	return nil
}
