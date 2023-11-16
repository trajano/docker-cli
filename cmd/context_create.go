/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:       "create CONTEXT TARGET",
	Args:      cobra.ExactArgs(2),
	ValidArgs: []string{"CONTEXT", "TARGET"},
	Example:   "docker context create nas ssh://nas",
	Short:     "Create a Docker context",
	Long: `Creates a new Docker context.  This effectively maps to:

  docker context create CONTEXT --docker host=TARGET
`,
	Run: func(cmd *cobra.Command, args []string) {
		RunDockerCommand("context", "create", args[0], "--docker", fmt.Sprintf("host=%s", args[1]))
	},
}

func init() {
	contextCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
