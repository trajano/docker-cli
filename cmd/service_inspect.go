/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
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

		var services []swarm.Service
		if err := json.Unmarshal(output, &services); err != nil {
			fmt.Println(err)
			return err
		}
		for i := range services {
			services[i].PreviousSpec = nil
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
