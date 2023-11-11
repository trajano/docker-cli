package main

import (
	"fmt"
	"os"
)

type Command interface {
	Process()
}

type BaseCommand struct {
	Receiver string
}

func (c BaseCommand) Process() {}

func ServiceCommandGroup() bool {
	args := os.Args[1:]
	if len(args) < 3 {
		return false
	}

	service := args[1]
	command := args[2]

	fmt.Printf("Service: %s, Command: %s\n", service, command)
	return true
}


func main() {
	args := os.Args[1:]

	switch {
	case len(args) == 0:
		fmt.Println("No command provided")
		os.Exit(1)
	case args[0] == "service" && ServiceCommandGroup():
		// Do nothing, as the service command is already handled in ServiceCommandGroup
	case args[0] == "ps":
		CmdPs(args...)
	default:
		RunDockerCommand(args...)
	}
}
