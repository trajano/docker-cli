package docker

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// Given a search key provide the services.
//
// Parameters:
// s - search key.  If it starts with `~` then this will do a partial name search within the names.  If it is nil all service are returned
//
// Returns:
// a list of service
func Services(keys []string) ([]swarm.Service, error) {
	var services []swarm.Service
	var err error
	_, isDockerHostSet := os.LookupEnv("DOCKER_HOST")
	if isDockerHostSet {
		services, err = servicesViaCli()
	} else {
		services, err = servicesViaClient()
	}
	if err != nil {
		return nil, err
	}
	var filteredServices []swarm.Service
	for _, service := range services {
		if len(keys) == 0 || isServiceSatisfiedBySearchKey(keys, &service) {
			filteredServices = append(filteredServices, service)
		}
	}
	return filteredServices, nil
}

// Given a search key provide the service names.
//
// Parameters:
// s - search key.  If it starts with `~` then this will do a partial name search within the names.  If it is nil all service are returned
//
// Returns:
// a list of service names
func ServiceNames(keys []string) ([]string, error) {
	var services []swarm.Service
	var err error
	_, isDockerHostSet := os.LookupEnv("DOCKER_HOST")
	if isDockerHostSet {
		services, err = servicesViaCli()
	} else {
		services, err = servicesViaClient()
	}
	if err != nil {
		return nil, err
	}
	var serviceNames []string
	for _, service := range services {
		if len(keys) == 0 || isServiceSatisfiedBySearchKey(keys, &service) {
			serviceNames = append(serviceNames, service.Spec.Name)
		}
	}
	return serviceNames, nil
}

type serviceLsInfo struct {
	ID       string `json:"ID"`
	Image    string `json:"Image"`
	Mode     string `json:"Mode"`
	Name     string `json:"Name"`
	Ports    string `json:"Ports"`
	Replicas string `json:"Replicas"`
}

// Workaround for when direct API access is not available due to using DOCKER_HOST and fallback to the CLI
func servicesViaCli() ([]swarm.Service, error) {
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
		var serviceLsInfo serviceLsInfo
		var runningCount int
		var desiredCount int
		if err := json.Unmarshal([]byte(line), &serviceLsInfo); err != nil {
			return nil, err
		}
		if runningCount, err = strconv.Atoi(strings.Split(serviceLsInfo.Replicas, "/")[0]); err != nil {
			return nil, err
		}
		if desiredCount, err = strconv.Atoi(strings.Split(serviceLsInfo.Replicas, "/")[1]); err != nil {
			return nil, err
		}

		service := swarm.Service{
			ID: serviceLsInfo.ID,
			Spec: swarm.ServiceSpec{
				Annotations: swarm.Annotations{
					Name: serviceLsInfo.Name,
				},
				TaskTemplate: swarm.TaskSpec{
					ContainerSpec: &swarm.ContainerSpec{
						Image: serviceLsInfo.Image,
					},
				},
			},
			ServiceStatus: &swarm.ServiceStatus{
				RunningTasks: uint64(runningCount),
				DesiredTasks: uint64(desiredCount),
			},
		}
		services = append(services, service)
	}
	return services, nil
}

func servicesViaClient() ([]swarm.Service, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return cli.ServiceList(context.Background(), types.ServiceListOptions{Status: true})
}

func isServiceSatisfiedBySearchKey(keys []string, service *swarm.Service) bool {
	for _, s := range keys {
		if (strings.HasPrefix(s, "~") && strings.Contains(service.Spec.Name, s[1:])) || service.Spec.Name == s || strings.HasPrefix(service.ID, s) {
			return true
		}
	}
	return false
}
