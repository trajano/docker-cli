/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var buildTag string
var buildQuiet bool
var buildPlain bool
var buildSecrets []string

/**
 * appends a secret to the slice if it exists
 */
func appendSecret(secrets []string, secretName string, pathElem ...string) ([]string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	secretPath := filepath.Join(append([]string{currentUser.HomeDir}, pathElem...)...)
	if _, err := os.Stat(secretPath); err == nil {
		return append(secrets, "id="+secretName+",src="+secretPath), nil
	} else {
		return secrets, nil
	}
}

var buildCmd = &cobra.Command{
	Use:   "build [PATH | URL | -]",
	Short: "Build an image from a Dockerfile",
	Long:  "Build an image from a Dockerfile, implicitly adds secrets",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {

		buildSecrets, err := appendSecret(buildSecrets, "init-gradle", ".gradle", "init.gradle")
		if err != nil {
			return err
		}
		buildSecrets, err = appendSecret(buildSecrets, "npmrc", ".npmrc")
		if err != nil {
			return err
		}
		buildSecrets, err = appendSecret(buildSecrets, "settings-xml", ".m2", "settings.xml")
		if err != nil {
			return err
		}
		buildSecrets, err = appendSecret(buildSecrets, "aws-credentials", ".aws", "credentials")
		if err != nil {
			return err
		}

		var flags []string
		flags = append(flags, "--load")
		for _, buildSecret := range buildSecrets {
			flags = append(flags, "--secret="+buildSecret)
		}

		if buildTag != "" {
			flags = append(flags, "-t", buildTag)
		}

		if buildPlain {
			flags = append(flags, "--progress=plain")
		}

		if buildQuiet {
			flags = append(flags, "--quiet")
		}

		if len(args) == 0 {
			RunDockerCommand(append(append([]string{"buildx", "build"}, flags...), ".")...)
		} else {
			RunDockerCommand(append(append([]string{"buildx", "build"}, flags...), args[0])...)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&buildTag, "tag", "t", "", "Name and optionally a tag (format: \"name:tag\")")
	buildCmd.Flags().StringArrayVar(&buildSecrets, "secret", []string{}, "Secret to expose to the build (format: \"id=mysecret[,src=/local/secret]\")")
	buildCmd.Flags().BoolVarP(&buildQuiet, "quiet", "q", false, "Suppress the build output and print image ID on success")
	buildCmd.Flags().BoolVarP(&buildPlain, "plain", "", false, "Plain output")
}
