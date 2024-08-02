/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system COMMAND",
	Short: "Manage Docker",
	Long:  `Manage Docker`,
}

func init() {
	rootCmd.AddCommand(systemCmd)
}
