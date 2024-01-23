/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var PsAll bool

// Sanitizes the image name.  This strips off `:latest` tag and `ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib` is replaced with `opentelemetry-collector`
// This also handles images containing a `library/` path  so those will be stripped off along with the prefix as it is assumed to be proxied.
func sanitizeImageName(input string) string {
  parts := strings.Split(input, "@sha256:")
	t := parts[0]

	if strings.HasSuffix(input, ":latest") {
		t = strings.TrimSuffix(input, ":latest")
	}
	if strings.HasPrefix(t, "ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib") {
		return "<opentelemetry-collector>"
	}
	if index := strings.Index(t, "/library/"); index != -1 {
		return t[index+9:]
	}

	return t
}

func psFunc(cmd *cobra.Command, keys []string) error {
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
	t.AppendHeader(table.Row{"Name", "Image", "", "Status", "Since", "", "Ports"})
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

				_, swarmService := container.Labels["com.docker.swarm.task"]

				imageName := sanitizeImageName(container.Image)
				if strings.HasPrefix(container.Image, "sha256:") {
					imageName = container.Image[7:19]
				}

				var exposedPorts []string
				ports := container.Ports
				if swarmService {
					ports = []types.Port{}
				}
				for _, port := range ports {
					if port.PublicPort != 0 {
						exposedPorts = append(exposedPorts, strconv.FormatUint(uint64(port.PublicPort), 10))
					}
				}
				t.AppendRow([]interface{}{
					container.Names[0][1:],
					imageName,
					healthState,
					containerJson.State.Status,
					startedAt.Format("2006-01-02 15:04"),
					humanize.Time(startedAt),
					strings.Join(exposedPorts, ", "),
				})
				// t.AppendRow([]interface{}{container.Names[0][1:], container.Image, containerJson.State.Health.Status})

			}
		}
	}
	t.Render()
	return nil

}

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:     "ps",
	Aliases: []string{"ls"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: psFunc,
}

func init() {
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolVarP(&PsAll, "all", "a", false, "Show all containers")
}
