/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var shCmd = &cobra.Command{
	Use:   "sh CONTAINERID",
	Short: "Executes a Bourne shell in a running container",
	Long:  `Executes a Bourne shell via /bin/sh in a running container`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		RunDockerCommand("exec", "-it", args[0], "/bin/sh")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(shCmd)
}
