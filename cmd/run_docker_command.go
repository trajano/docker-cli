package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunDockerCommand(args ...string) {
	dockerCommand := append([]string{"docker"}, args...)
	if displayRunCommand {
		fmt.Println(strings.Join(dockerCommand, " "))
		os.Exit(0)
	}

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
}

// Run a docker command and extract individual lines as bytes suitable for json Unmarshal
func RunDockerCommandMapJsonBytes(args ...string) ([][]byte, error) {
	dockerCommand := append([]string{"docker"}, args...)
	cmd := exec.Command(dockerCommand[0], dockerCommand[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			return nil, err
		}
	}
	lines := strings.Split(string(output), "\n")
	var collect [][]byte
	for _, line := range lines {
		if len(line) > 0 {
			collect = append(collect, []byte(line))
		}
	}
	return collect, nil
}

type DockerContainer struct {
	ID    string   `json:"ID"`
	Names []string `json:"Names"`
}

func RunDockerPs() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
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
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
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
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
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
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}
	var containerIDs []string
	for _, container := range containers {
		containerIDs = append(containerIDs, container.ID)
	}
	return containerIDs, nil
}
