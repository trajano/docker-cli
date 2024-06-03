/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var composeCmd = &cobra.Command{
	Use:   "compose [OPTIONS] COMMAND",
	Short: "Define and run multi-container applications with Docker",
	Long:  `Define and run multi-container applications with Docker`,
	Run: func(cmd *cobra.Command, args []string) {
		RunDockerCommand(append([]string{"compose"}, args...)...)
	},
}

func init() {
	rootCmd.AddCommand(composeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
