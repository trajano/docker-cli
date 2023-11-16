/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"trajano.net/docker-cli/docker"
)

var serviceRestartDetached bool

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a service",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, keys []string) error {
		services, err := docker.ServiceNames(keys)
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, serviceName := range services {
			if serviceRestartDetached {
				RunDockerCommand("service", "update", "-d", "--force", serviceName)
			} else {
				RunDockerCommand("service", "update", "--force", serviceName)
			}
		}
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(restartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	restartCmd.Flags().BoolVarP(&serviceRestartDetached, "detached", "d", false, "Restart in background")
}
