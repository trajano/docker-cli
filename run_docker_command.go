package main

import (
	"fmt"
	"os"
	"os/exec"
)

func RunDockerCommand(args ...string) {
	dockerCommand := append([]string{"docker"}, args...)
	cmd := exec.Command(dockerCommand[0], dockerCommand[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			fmt.Printf("Error running Docker command: %v\n", err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
