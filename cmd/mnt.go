/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var MntVolume string
var MntImage string

var mntCmd = &cobra.Command{
	Use:   "mnt VOLUME|PATH [IMAGE] [ARGS...]",
	Short: "Mounts a volume or path and runs a command",
	Long:  `Mounts a volume or path (may be relative) then runs the specified Docker image with logging driver turned off and allows interaction from the console`,
    DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		volume := "."
		if len(args) >= 1 {
			volume = args[0]
		}
		// ignore the error just presume it's a volume label in that case
		stat, err := os.Stat(volume)
		if err == nil && stat.IsDir() && (strings.HasPrefix(volume, "/") ||
			strings.HasPrefix(volume, "./") ||
			strings.HasPrefix(volume, "../") ||
			volume == "." ||
			volume == "..") {
			path, err := filepath.Abs(volume)
			if err != nil {
				return err
			}
			volume = path
		}
		image := "alpine"
		var cmdArgs []string
		if len(args) >= 2 {
			image = args[1]
		}
		if len(args) > 2 {
			cmdArgs = args[2:]
		}

		RunDockerCommand(append([]string{
			"run",
			"-it",
			"-v", volume + ":/mnt",
			"--log-driver", "none",
			"--workdir", "/mnt",
			"--rm",
			image},
			cmdArgs...)...)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mntCmd)

}
