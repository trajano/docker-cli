/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var bashCmd = &cobra.Command{
	Use:   "bash CONTAINERID",
	Short: "Executes a bash shell in a running container",
	Long:  `Executes a bash shell via /bin/bash in a running container`,
	Args:  cobra.ExactArgs(1),
  SilenceErrors: false,
  SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		RunDockerCommand("exec", "-it", args[0], "/bin/bash")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(bashCmd)
}
