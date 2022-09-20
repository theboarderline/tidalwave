/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"tidalwave/internal/google"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

// controlplaneCmd represents the controlplane command
var controlplaneCmd = &cobra.Command{
	Use:   "controlplane",
	Short: "Create/Modify/Delete a DevOps controlplane",
	Long:  "Create/Modify/Delete a DevOps controlplane",
}

func init() {
	rootCmd.AddCommand(controlplaneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// controlplaneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// controlplaneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CreateGoogleControlplane() (*google.Controlplane, error) {
	viper.SetDefault("spec.region", "us-central1")
	viper.SetDefault("spec.cidrs.nodes", "10.0.0.0/24")
	viper.SetDefault("spec.cidrs.pods", "10.1.0.0/16")
	viper.SetDefault("spec.cidrs.services", "10.2.0.0/20")
	viper.SetDefault("spec.cluster.machineType", "n2-standard-4")
	viper.SetDefault("spec.cluster.minNodeCount", 1)
	viper.SetDefault("spec.cluster.maxNodeCount", 3)
	viper.SetDefault("spec.cluster.masterAuthBlock", []map[string]string{
		{
			"DisplayName": "public",
			"CidrBlock":   "0.0.0.0/0",
		},
	})
	viper.SetDefault("spec.cluster.masterCidrBlock", "172.16.0.0/28")
	name := viper.GetString("metadata.name")
	projectId := viper.GetString("spec.projectId")
	region := viper.GetString("spec.region")
	nodesCidr := viper.GetString("spec.cidrs.nodes")
	podCidr := viper.GetString("spec.cidrs.pods")
	serviceCidr := viper.GetString("spec.cidrs.services")
	machineType := viper.GetString("spec.cluster.machineType")
	diskSize := viper.GetInt32("spec.cluster.diskSize")
	minNodes := viper.GetInt32("spec.cluster.minNodeCount")
	maxNodes := viper.GetInt32("spec.cluster.maxNodeCount")
	masterAuthCidrBlocks := []*containerpb.MasterAuthorizedNetworksConfig_CidrBlock{}
	if err := viper.UnmarshalKey("spec.cluster.masterAuthBlock", &masterAuthCidrBlocks); err != nil {
		return nil, err
	}
	masterIpv4CidrBlock := viper.GetString("spec.cluster.masterCidrBlock")
	cp := google.Controlplane{
		Vpc: google.Vpc{
			Name:      name,
			ProjectID: projectId,
		},
		Subnetwork: google.Subnetwork{
			Name:         fmt.Sprintf("%s-controlplane", name),
			ProjectID:    projectId,
			Region:       region,
			NodesCidr:    nodesCidr,
			PodsCidr:     podCidr,
			ServicesCidr: serviceCidr,
		},
		Router: google.Router{
			Name:      name,
			ProjectID: projectId,
			Region:    region,
		},
		Cluster: google.Cluster{
			Name:                 fmt.Sprintf("%s-controlplane", name),
			ProjectID:            projectId,
			Region:               region,
			Network:              name,
			Subnetwork:           fmt.Sprintf("%s-controlplane", name),
			MachineType:          machineType,
			DiskSizeGb:           diskSize,
			MinNodeCount:         minNodes,
			MaxNodeCount:         maxNodes,
			MasterAuthCidrBlocks: masterAuthCidrBlocks,
			MasterIpv4CidrBlock:  masterIpv4CidrBlock,
		},
	}
	return &cp, nil
}
