package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cristalhq/acmd"
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
	cmds := []acmd.Command{
		{
			Name:        "ps",
			Description: "docker ps with better formatting",
			ExecFunc: func(ctx context.Context, args []string) error {
				RunDockerCommand("ps", "--format", "table {{.Names}}\t{{.Image}}\t{{.Status}}")
				return nil
			},
		},
		{
			Name:        "run",
			Description: "docker ps with better formatting",
			ExecFunc: func(ctx context.Context, args []string) error {
				RunDockerCommand(append([]string{"run", "-it", "--log-driver=none", "--rm"}, args...)...)
				return nil
			},
		},
		{
			Name:        "ptag",
			Description: "Tag and Push",
			ExecFunc: func(ctx context.Context, args []string) error {
				RunDockerCommand("tag", args[0])
				RunDockerCommand("push", args[1])
				return nil
			},
		},
		{
			Name:        "service",
			Description: "Tag and Push",
			Subcommands: []acmd.Command{
				{Name: "restart",
					Description: "Restart a service",
					ExecFunc: func(ctx context.Context, args []string) error {
						RunDockerCommand("service", "update", "--force", args[0])
						return nil
					},
				},
			},
		},
	}
	r := acmd.RunnerOf(cmds, acmd.Config{})
	if err := r.Run(); err != nil {
		RunDockerCommand(os.Args[1:]...)
		//r.Exit(err)
	}
	fmt.Println("Bar")
}

// github.com/cristalhq/acmd	parser := argparse.NewParser("docker-cli", "Wraps the Docker command")
// 	var firstArgument = parser.StringPositional(&argparse.Options{Required: true})
// 	err := parser.Parse(os.Args)
// 	if err != nil {
// 		fmt.Print(parser.Usage(err))
// 		RunDockerCommand(os.Args[1:]...)
// 	} else {
// 		fmt.Println("something" + *firstArgument)
// 	}
// 	// 	args := os.Args[1:]

// 	// // load up the commands here

// // switch {
// // case len(args) == 0:
// //
// //	fmt.Println("No command provided")
// //	os.Exit(1)
// //
// // case args[0] == "service" && ServiceCommandGroup():
// //
// //	// Do nothing, as the service command is already handled in ServiceCommandGroup
// //
// // case args[0] == "ps":
// //
// //	CmdPs(args...)
// //
// // default:
// //
// //		RunDockerCommand(args...)
// //	}
//}
