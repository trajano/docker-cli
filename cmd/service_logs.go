/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"time"

	"github.com/spf13/cobra"
)

var serviceLogsSince string
var serviceLogsFollow bool
var serviceLogsSincePositional bool = false

// serviceLogsCmd represents the logs command
var serviceLogsCmd = &cobra.Command{
	Use: "logs [SINCE] SERVICE|TASK",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 2 {
			_, err := time.Parse(time.RFC3339Nano, args[0])
			if err != nil {
				_, err = time.ParseDuration(args[0])
			}
			if err != nil {
				return err
			}
			serviceLogsSince = args[0]
			serviceLogsSincePositional = true
			return nil
		} else if len(args) == 1 {
			return nil
		} else {
			return errors.New("invalid number of arguments, expected 1 or 2")
		}
	},
	Short: "Fetch the logs of a service or task",
	Long:  `Fetch the logs of a service or task.  This defaults to start following from the end of the log.`,
	Run: func(cmd *cobra.Command, args []string) {
		dockerCmd := []string{"service", "logs", "--raw", "--since", serviceLogsSince}
		if serviceLogsFollow {
			dockerCmd = append(dockerCmd, "--follow")
		}
		if serviceLogsSincePositional {
			dockerCmd = append(dockerCmd, args[1])
		} else {
			dockerCmd = append(dockerCmd, args...)
		}
		RunDockerCommand(dockerCmd...)
	},
}

func init() {
	serviceCmd.AddCommand(serviceLogsCmd)
	serviceLogsCmd.Flags().StringVarP(&serviceLogsSince, "since", "", "0s", "Show logs since timestamp (e.g, \"2013-01-02T13:23:37Z\") or relative (e.g. \"42m\" for 42 minutes)")
	serviceLogsCmd.Flags().BoolVarP(&serviceLogsFollow, "follow", "f", true, "Follow log output")
}
