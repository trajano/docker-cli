/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [CONTAINERID...]",
	Short: "Remove a container",
	Long:  `Kill and remove a container along with any anonymous volumes associated with it and prune`,
	RunE: func(cmd *cobra.Command, containerIds []string) error {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return err
		}
		for _, containerId := range containerIds {
			err := cli.ContainerRemove(context.Background(), containerId, container.RemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			})
			if err != nil {
				return err
			}
			fmt.Println(containerId)
		}
		containerPruneReport, err := cli.ContainersPrune(context.Background(), filters.Args{})
		for _, prunedContainerId := range containerPruneReport.ContainersDeleted {
			fmt.Printf("container %s pruned\n", prunedContainerId)
		}
		if err != nil {
			return err
		}
		volumesPruneReport, err := cli.VolumesPrune(context.Background(), filters.Args{})
		for _, prunedVolumeId := range volumesPruneReport.VolumesDeleted {
			fmt.Printf("volume %s pruned\n", prunedVolumeId)
		}
		if err != nil {
			return err
		}
		if containerPruneReport.SpaceReclaimed > 0 {
			fmt.Printf("%s freed from containers\n", humanize.Bytes(containerPruneReport.SpaceReclaimed))
		}
		if volumesPruneReport.SpaceReclaimed > 0 {
			fmt.Printf("%s freed from volumes\n", humanize.Bytes(volumesPruneReport.SpaceReclaimed))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
