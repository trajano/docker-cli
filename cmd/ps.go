/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var PsAll bool

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, keys []string) error {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return err
		}
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
			All: PsAll,
		})

		if err != nil {
			return err
		}
		t := table.NewWriter()
		t.SetStyle(table.StyleColoredDark)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Image", "", "Status", "Since"})
		for _, container := range containers {
			if len(keys) == 0 || IsContainerSatisfiedBySearchKey(keys, &container) {
				if containerJson, err := cli.ContainerInspect(context.Background(), container.ID); err != nil {
					return err
				} else {
					var healthState = ""
					if containerJson.State.Health != nil {
						healthState = " "
						if containerJson.State.Health.Status == "healthy" {
							healthState = ""
						} else if containerJson.State.Health.Status == "unhealthy" {
							healthState = ""
						} else if containerJson.State.Health.Status == "starting" {
							healthState = ""
						}
					}
					// createdAt, err := time.Parse(time.RFC3339Nano, containerJson.Created)
					// if err != nil {
					// 	return err
					// }
					startedAt, err := time.Parse(time.RFC3339Nano, containerJson.State.StartedAt)
					if err != nil {
						return err
					}
					// finishedAt, err := time.Parse(time.RFC3339Nano, containerJson.State.FinishedAt)
					// if err != nil {
					// 	return err
					// }

					// var startupTime time.Duration
					// if !startedAt.IsZero() && !finishedAt.IsZero() {
					// 	startupTime = startedAt.Sub(finishedAt)
					// } else if !startedAt.IsZero() && finishedAt.IsZero() {
					// 	startupTime = startedAt.Sub(createdAt)
					// }

					imageName := container.Image
					if strings.HasPrefix(container.Image, "sha256:") {
						imageName = container.Image[7:19]
					}

					t.AppendRow([]interface{}{
						container.Names[0][1:],
						imageName,
						healthState,
						containerJson.State.Status,
						startedAt.Format("2006-01-02 15:04"),
						humanize.Time(startedAt),
					})
					// t.AppendRow([]interface{}{container.Names[0][1:], container.Image, containerJson.State.Health.Status})

				}
			}
		}
		t.Render()
		return nil

	},
}

func init() {
	rootCmd.AddCommand(psCmd)
	// psCmd.SilenceUsage = true
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// psCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	psCmd.Flags().BoolVarP(&PsAll, "all", "a", false, "Show all containers")
}
