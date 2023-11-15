package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jedib0t/go-pretty/v6/table"
)

func RunDockerCommand(args ...string) {
	dockerCommand := append([]string{"docker"}, args...)
	cmd := exec.Command(dockerCommand[0], dockerCommand[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			fmt.Printf("Error running Docker command: %v\n", err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

type DockerContainer struct {
	ID    string   `json:"ID"`
	Names []string `json:"Names"`
}

func RunDockerPs() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		// fmt.Println("Error listing containers:", err)
		panic("Error running docker ps")
	}
	for _, container := range containers {
		fmt.Println(container.Names)
		fmt.Println(container.ID)
	}
	// cmd := exec.Command("docker", "ps", "--format", "json", "--no-trunc")
	// output, err := cmd.Output()
	// if err != nil {
	//   panic("Error running docker ps")
	// }
	// cmd.Stdout = os.Stdout

}

// Given a search key provide the container IDs
//
// Parameters:
// s - search key.  If it starts with `~` then this will do a partial name search within the names.
//
// Returns:
// a list of container IDs
func Containers(s string) ([]string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	var containerIDs []string
	for _, container := range containers {
		for _, name := range container.Names {
			if (strings.HasPrefix(s, "~") && strings.Contains(name[1:], s[1:])) || name[1:] == s[1:] {
				containerIDs = append(containerIDs, container.ID)
				break // Break the inner loop as the container ID is found
			}
		}
	}
	return containerIDs, nil
}

// Given a search key provide the container IDs
//
// Parameters:
// keys - search keys.  If it starts with `~` then this will do a partial name search within the names.  If the list is empty, this will return all container IDs.
//
// Returns:
// a list of container IDs
func Containers2(keys []string) ([]string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	var containerIDs []string
	for _, container := range containers {
		if len(keys) == 0 || IsContainerSatisfiedBySearchKey(keys, &container) {
			containerIDs = append(containerIDs, container.ID)
		}
	}
	return containerIDs, nil
}

func PrettyPs(keys []string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		return err
	}
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredDark)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Image", "", "Status", "Since", "Startup"})
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
				// startedAt, err := time.Parse(time.RFC3339Nano, containerJson.State.StartedAt)
				// if err != nil {
				// 	return err
				// }
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
				t.AppendRow([]interface{}{container.Names[0][1:], container.Image, healthState, containerJson.State.Status, container.Status})
				// t.AppendRow([]interface{}{container.Names[0][1:], container.Image, containerJson.State.Health.Status})

			}
		}
	}
	t.Render()
	return nil
}

func IsContainerSatisfiedBySearchKey(keys []string, container *types.Container) bool {
	for _, s := range keys {
		for _, name := range container.Names {
			if (strings.HasPrefix(s, "~") && strings.Contains(name[1:], s[1:])) || name[1:] == s || strings.HasPrefix(container.ID, s) {
				return true
			}
		}
	}
	return false
}

// Given a search key provide the container IDs
//
// Parameters:
// s - search key.  If it starts with `~` then this will do a partial name search within the names.  If it is nil all container IDs are passed
//
// Returns:
// a list of container IDs
func AllContainers() ([]string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	var containerIDs []string
	for _, container := range containers {
		containerIDs = append(containerIDs, container.ID)
	}
	return containerIDs, nil
}
