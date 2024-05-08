/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// duCmd represents the du command
var duCmd = &cobra.Command{
	Use:   "du",
	Short: "Disk usage report",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return err
		}
		diskUsage, err := cli.DiskUsage(context.Background(), types.DiskUsageOptions{
			Types: []types.DiskUsageObject{
				types.ContainerObject,
				types.BuildCacheObject,
				types.VolumeObject,
				types.ImageObject,
			},
		})
		if err != nil {
			return err
		}
		jsonData, err := json.MarshalIndent(diskUsage, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(duCmd)
}
