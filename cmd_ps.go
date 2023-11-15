package main

func CmdPs(args ...string) {
	RunDockerCommand(append([]string{"ps", "--format", "table {{.Names}}\t{{.Image}}\t{{.Status}}"}, args[1:]...)...)
}
