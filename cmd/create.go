/*
Package cmd is the entrypoint the for cli
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"tidalwave/internal/tidalwave"

	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DevOps controlplane cluster",
	Long:  "Create a DevOps controlplane cluster",
	Run: func(cmd *cobra.Command, args []string) {
		switch viper.Get("spec.provider") {
		case "google":
			emoji.Println(":joystick: Create Google Controlplane")
			c, err := CreateGoogleControlplane()
			if err != nil {
				log.Fatal(err)
			}
			if err := tidalwave.CheckApis(c); err != nil {
				log.Fatal(err)
			}
			err = tidalwave.CreateCluster(c)
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
