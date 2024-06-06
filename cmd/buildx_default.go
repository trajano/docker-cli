/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var buildxDefaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Sets the builder to the default",
	Long:  `Sets the builder to the default`,
	Run: func(cmd *cobra.Command, args []string) {
		RunDockerCommand("buildx", "use", "default")
	},
}

func init() {
	buildxCmd.AddCommand(buildxDefaultCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
