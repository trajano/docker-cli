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
	"gopkg.in/yaml.v3"
	"trajano.net/docker-cli/docker"
)

type serviceEnvInfo struct {
	Name string
	Env  []string
}

// serviceEnvCmd represents the inspect command
var serviceEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Display environment settings for one or more services",
	Long:  `Display environment settings for one or more services.  The output will be in YAML as it is likely to be copied and pasted back into a docker-compose.yml file.`,
	RunE: func(cmd *cobra.Command, keys []string) error {

		var serviceNames []string
		var err error
		if serviceNames, err = docker.ServiceNames(keys); err != nil {
			return err
		}
		dockerCommandArgs := append([]string{"docker", "service", "inspect", "--format", "json"}, serviceNames...)
		dockerCommand := exec.Command(dockerCommandArgs[0], dockerCommandArgs[1:]...)
		dockerCommand.Stdin = os.Stdin
		dockerCommand.Stderr = os.Stderr
		output, err := dockerCommand.Output()

		var envInfos []serviceEnvInfo
		var services []swarm.Service
		if err := json.Unmarshal(output, &services); err != nil {
			fmt.Println(err)
			return err
		}
		for _, service := range services {
			var envInfo serviceEnvInfo
			envInfo.Name = service.Spec.Name
			envInfo.Env = service.Spec.TaskTemplate.ContainerSpec.Env
			envInfos = append(envInfos, envInfo)
		}
		var marshalledBytes []byte
		if len(services) == 1 {
			marshalledBytes, err = yaml.Marshal(envInfos[0].Env)
		} else {
			marshalledBytes, err = yaml.Marshal(envInfos)
		}
		if err != nil {
			return err
		}
		fmt.Println(string(marshalledBytes))
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(serviceEnvCmd)
}
