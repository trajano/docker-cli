package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type CommandReceiver interface {
	runDockerCommand(args ...string)
}
type Command interface {
	execute(r *CommandReceiver)
}

type Invoker struct {
	command Command
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
	parser := argparse.NewParser("docker-cli", "Wraps the Docker command")
	var firstArgument = parser.StringPositional(&argparse.Options{Required: true})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		RunDockerCommand(os.Args[1:]...)
	} else {
		fmt.Println("something" + *firstArgument)
	}
	// 	args := os.Args[1:]

	// // load up the commands here

	// switch {
	// case len(args) == 0:
	//
	//	fmt.Println("No command provided")
	//	os.Exit(1)
	//
	// case args[0] == "service" && ServiceCommandGroup():
	//
	//	// Do nothing, as the service command is already handled in ServiceCommandGroup
	//
	// case args[0] == "ps":
	//
	//	CmdPs(args...)
	//
	// default:
	//
	//		RunDockerCommand(args...)
	//	}
}
