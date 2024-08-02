/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var systemDialStdioCmd = &cobra.Command{
	Use:   "dial-stdio",
	Short: "Manage Docker",
	Args:  cobra.NoArgs,
	Long:  `Proxy the stdio stream to the daemon connection. Should not be invoked manually.`,
	Run: func(cmd *cobra.Command, services []string) {
		RunDockerCommand("system", "dial-stdio")
	},
}

func init() {
	systemCmd.AddCommand(systemDialStdioCmd)
}
