/*
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

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a DevOps controlplane cluster",
	Long:  "Update a DevOps controlplane cluster",
	Run: func(cmd *cobra.Command, args []string) {
		switch viper.Get("spec.provider") {
		case "google":
			emoji.Println("Update Google Controlplane :joystick:")
			c, err := CreateGoogleControlplane()
			if err != nil {
				log.Fatal(err)
			}
			err = tidalwave.UpdateCluster(c)
			if err != nil {
				log.Fatal(err)
			}
		case "aws":
			fmt.Println("Configure AWS controlplane")
		}
	},
}

func init() {
	controlplaneCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
