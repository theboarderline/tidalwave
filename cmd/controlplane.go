/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
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
