/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types/swarm"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"trajano.net/docker-cli/docker"
)

var lsDown bool

// restartCmd represents the restart command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Restart a service",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, keys []string) error {
		t := table.NewWriter()
		t.SetStyle(table.StyleColoredDark)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Image", "Replicas"})
		var services []swarm.Service
		var err error
		if services, err = docker.Services(keys); err != nil {
			return err
		}
		for _, service := range services {

			replicas := fmt.Sprintf("%d/%d", service.ServiceStatus.RunningTasks, service.ServiceStatus.DesiredTasks)
			if (lsDown && service.ServiceStatus.DesiredTasks > 0 && service.ServiceStatus.RunningTasks != service.ServiceStatus.DesiredTasks) || (!lsDown) {
				t.AppendRow([]interface{}{
					service.Spec.Name,
					stripLatestSuffix(service.Spec.TaskTemplate.ContainerSpec.Image),
					replicas,
				})
			}
		}
		t.Render()
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	lsCmd.Flags().BoolVarP(&lsDown, "down", "d", false, "Show only services that are down")
}
