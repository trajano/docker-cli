/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var RunDetach bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an image",
	Long:  `Runs the specified Docker image with logging driver turned off and allows interaction from the console`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if RunDetach {
			RunDockerCommand(append([]string{"run", "--detach", "--rm"}, args...)...)
		} else {

			RunDockerCommand(append([]string{"run", "-it", "--log-driver=none", "--rm"}, args...)...)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().BoolVarP(&RunDetach, "detach", "d", false, "Run the container in the background and print the container ID")
}
