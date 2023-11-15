/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
)

var patch bool

// ptagCmd represents the ptag command
var ptagCmd = &cobra.Command{
	Use:   "ptag",
	Short: "Tag and Push",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if patch {
			patched := createPatchImageName(args[0])
			RunDockerCommand("tag", args[0], patched)
			fmt.Printf("docker tag %s %s\n", args[0], patched)
			RunDockerCommand("push", patched)
		} else {
			RunDockerCommand("tag", args[0], args[1])
			RunDockerCommand("push", args[1])
		}
	},
}

func createPatchImageName(input string) string {
	re := regexp.MustCompile(`^(?P<repository>[^:]+):(?P<tag>.+)$`)
	match := re.FindStringSubmatch(input)

	if len(match) > 0 {
		repository := match[1]
		patchTag := "patch-" + strconv.Itoa(rand.Intn(10000))
		return repository + ":" + patchTag
	}

	return input
}
func init() {
	rootCmd.AddCommand(ptagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ptagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	ptagCmd.Flags().BoolVarP(&patch, "patch", "p", false, "Generate patch label")
}
