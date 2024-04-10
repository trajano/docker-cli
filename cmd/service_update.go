/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var updateImage string

// restartCmd represents the restart command
var serviceUpdateCmd = &cobra.Command{
	Use:   "update",
	Args:  cobra.ArbitraryArgs,
	Short: "Update a service",
	Long:  `Update a service`,
	Run: func(cmd *cobra.Command, services []string) {

    if updateImage != "" {
      for _, service := range services {
        RunDockerCommand("service", "update", "--with-registry-auth", "--image", updateImage, service)
      }
    } else {
      RunDockerCommand(append([]string{"service", "update"}, services...)...)
    }
	},
}

func init() {
	serviceCmd.AddCommand(serviceUpdateCmd)
	serviceUpdateCmd.Flags().StringVarP(&updateImage, "image", "i", "", "Service image tag. Implies with-registry-auth.")
}
