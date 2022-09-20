/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"tidalwave/internal/google"
	"tidalwave/internal/tidalwave"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("controlplane called")
		switch viper.Get("spec.provider") {
		case "google":
			fmt.Println("Configure Google controlplane")
			viper.SetDefault("spec.region", "us-central1")
			viper.SetDefault("spec.cidrs.nodes", "10.0.0.0/24")
			c := google.Controlplane{
				Name:      viper.GetString("metadata.name"),
				ProjectID: viper.GetString("spec.projectId"),
				Region:    viper.GetString("spec.region"),
				Subnetwork: google.Subnetwork{
					NodesCidr: viper.GetString("spec.cidrs.nodes"),
				},
			}
			err := tidalwave.CreateCluster(&c)
			if err != nil {
				log.Fatal(err)
			}
		case "aws":
			fmt.Println("Configure AWS controlplane")
		}
	},
}

func init() {
	controlplaneCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
