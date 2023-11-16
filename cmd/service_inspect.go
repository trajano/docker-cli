/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types/swarm"
	"github.com/spf13/cobra"
	"trajano.net/docker-cli/docker"
)

type networkInfo struct {
	ID     string
	Name   string
	Driver string
}

type serviceSwarmLimit struct {
	*swarm.Limit
}
type serviceSwarmResourceRequirements struct {
	*swarm.ResourceRequirements
	Limits *serviceSwarmLimit
}
type serviceSwarmTaskSpec struct {
	*swarm.TaskSpec
	Resources *serviceSwarmResourceRequirements
}
type serviceSwarmSpec struct {
	*swarm.ServiceSpec
	TaskTemplate serviceSwarmTaskSpec
}
type serviceInfo struct {
	*swarm.Service
	Spec     *serviceSwarmSpec
	Endpoint *swarm.Endpoint
}

func formatBytes(in int64) string {
	bytes := float64(in)
	suffixes := []string{"B", "K", "M", "G", "T", "P", "E", "Z", "Y"}
	base := 1024.0
	if bytes < base {
		return fmt.Sprintf("%.0f", bytes)
	}
	exp := int(math.Log(bytes) / math.Log(base))
	index := int(math.Min(float64(exp), float64(len(suffixes)-1)))
	value := bytes / math.Pow(base, float64(exp))
	return fmt.Sprintf("%.2f%s", value, suffixes[index])
}

func (limit *serviceSwarmLimit) MarshalJSON() ([]byte, error) {
	limitMap := make(map[string]interface{})
	if limit.NanoCPUs != 0 {
		limitMap["NanoCPUs"] = limit.NanoCPUs
	}
	if limit.MemoryBytes != 0 {
		limitMap["MemoryBytes"] = formatBytes(limit.MemoryBytes)
	}
	if limit.Pids != 0 {
		limitMap["Pids"] = limit.Pids
	}
	return json.Marshal(limitMap)
}

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Display detailed information on one or more services",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, keys []string) error {

		var serviceNames []string
		var jsonBytelines [][]byte
		var err error
		if serviceNames, err = docker.ServiceNames(keys); err != nil {
			return err
		}
		if jsonBytelines, err = RunDockerCommandMapJsonBytes("network", "ls", "--format", "json", "--no-trunc"); err != nil {
			return err
		}
		networks := make(map[string]networkInfo)
		for _, jsonBytes := range jsonBytelines {
			var network networkInfo
			if err := json.Unmarshal(jsonBytes, &network); err != nil {
				return err
			}
			networks[network.ID] = network
		}
		dockerCommandArgs := append([]string{"docker", "service", "inspect", "--format", "json"}, serviceNames...)
		dockerCommand := exec.Command(dockerCommandArgs[0], dockerCommandArgs[1:]...)
		dockerCommand.Stdin = os.Stdin
		dockerCommand.Stderr = os.Stderr
		output, err := dockerCommand.Output()

		var services []serviceInfo
		if err := json.Unmarshal(output, &services); err != nil {
			fmt.Println(err)
			return err
		}
		for i := range services {
			services[i].PreviousSpec = nil
			services[i].Spec.TaskTemplate.Placement.Platforms = nil
			for j := range services[i].Spec.TaskTemplate.Networks {
				services[i].Spec.TaskTemplate.Networks[j].Target = networks[services[i].Spec.TaskTemplate.Networks[j].Target].Name
			}
			for j := range services[i].Endpoint.VirtualIPs {
				services[i].Endpoint.VirtualIPs[j].NetworkID = networks[services[i].Endpoint.VirtualIPs[j].NetworkID].Name
			}
		}

		jsonData, err := json.Marshal(services)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(inspectCmd)
}
