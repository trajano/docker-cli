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
var serviceLogsSincePositional bool = false

// serviceLogsCmd represents the logs command
var serviceLogsCmd = &cobra.Command{
	Use: "logs",
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
		if serviceLogsSincePositional {
			RunDockerCommand(append([]string{"service", "logs", "--raw", "--since", serviceLogsSince, "--follow"}, args[1])...)
		} else {
			RunDockerCommand(append([]string{"service", "logs", "--raw", "--since", serviceLogsSince, "--follow"}, args...)...)
		}

	},
}

func init() {
	serviceCmd.AddCommand(serviceLogsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serviceLogsCmd.Flags().StringVarP(&serviceLogsSince, "since", "", "0s", "Show logs since timestamp (e.g, \"2013-01-02T13:23:37Z\") or relative (e.g. \"42m\" for 42 minutes)")
}
