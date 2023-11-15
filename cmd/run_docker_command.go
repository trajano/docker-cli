package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
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

func Services(keys []string) ([]string, error) {
	var services []swarm.Service
	var err error
	_, isDockerHostSet := os.LookupEnv("DOCKER_HOST")
	if isDockerHostSet {
		services, err = ServicesViaCli()
	} else {
		services, err = ServicesViaClient()
	}
	if err != nil {
		return nil, err
	}
	var serviceNames []string
	for _, service := range services {
		if len(keys) == 0 || IsServiceSatisfiedBySearchKey(keys, &service) {
			serviceNames = append(serviceNames, service.Spec.Name)
		}
	}
	return serviceNames, nil
}

type ServiceLsInfo struct {
	ID       string `json:"ID"`
	Image    string `json:"Image"`
	Mode     string `json:"Mode"`
	Name     string `json:"Name"`
	Ports    string `json:"Ports"`
	Replicas string `json:"Replicas"`
}

// Workaround for when direct API access is not available due to using DOCKER_HOST and fallback to the CLI
func ServicesViaCli() ([]swarm.Service, error) {
	cmd := exec.Command("docker", "service", "ls", "--format", "json")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	// Yes I know this is not a good way of doing it, still learning Go, I presume there's some way to do it via a stream rather than reading it all in memory
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var services []swarm.Service
	for _, line := range lines {
		var serviceLsInfo ServiceLsInfo
		if err := json.Unmarshal([]byte(line), &serviceLsInfo); err != nil {
			return nil, err
		}
		service := swarm.Service{
			ID: serviceLsInfo.ID,
			Spec: swarm.ServiceSpec{
				Annotations: swarm.Annotations{
					Name: serviceLsInfo.Name,
				},
			},
		}
		services = append(services, service)
	}
	return services, nil
}

func ServicesViaClient() ([]swarm.Service, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return cli.ServiceList(context.Background(), types.ServiceListOptions{Status: true})
}

func IsServiceSatisfiedBySearchKey(keys []string, service *swarm.Service) bool {
	for _, s := range keys {
		if (strings.HasPrefix(s, "~") && strings.Contains(service.Spec.Name, s[1:])) || service.Spec.Name == s || strings.HasPrefix(service.ID, s) {
			return true
		}
	}
	return false
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
