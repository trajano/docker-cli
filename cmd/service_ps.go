/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"trajano.net/docker-cli/docker"
)

var psAll bool
var psDown bool

type servicePsInfo struct {
	ID           string `json:"ID"`
	Image        string `json:"Image"`
	CurrentState string `json:"CurrentState"`
	DesiredState string `json:"DesiredState"`
	Name         string `json:"Name"`
	Error        string `json:"Error"`
	Node         string `json:"Node"`
	Ports        string `json:"Ports"`
	Indent       bool
}

func removeNonNumericEndings(input string) string {
	index := strings.Index(input, ".")
	_, err := strconv.Atoi(input[index+1:])
	if err != nil {
		return input[0:index]
	} else {
		return input
	}

}

func sanitizeError(input string) string {
	t := strings.TrimSuffix(strings.TrimPrefix(input, "\""), "\"")
	if t == "task: non-zero exit (137): dockerexec: unhealthy container" {
		return " 137"
	}
	return t
}

func sanitizeCurrentState(input string) string {
	t := strings.TrimSuffix(strings.TrimPrefix(input, "Running "), " ago")
	return t
}

// restartCmd represents the restart command
var servicePsCmd = &cobra.Command{
	Use:   "ps",
	Args:  cobra.ArbitraryArgs,
	Short: "List the tasks of one or more services",
	Long:  `List the tasks of one or more services`,
	RunE: func(cmd *cobra.Command, keys []string) error {

		t := table.NewWriter()
		t.SetStyle(table.StyleColoredDark)
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "ID", "Image", "Node", "", "Current State", "Error"})

		var serviceNames []string
		var jsonBytelines [][]byte
		var err error
		if serviceNames, err = docker.ServiceNames(keys); err != nil {
			return err
		}
		if jsonBytelines, err = RunDockerCommandMapJsonBytes(append([]string{"service", "ps", "--format", "json", "--no-trunc"}, serviceNames...)...); err != nil {
			return err
		}
		var lastName string
		primaryRunning := false
		for _, jsonBytes := range jsonBytelines {
			var task servicePsInfo
			if err := json.Unmarshal(jsonBytes, &task); err != nil {
				return err
			}
			task.Indent = lastName == task.Name
			lastName = task.Name
			name := removeNonNumericEndings(task.Name)
			if !task.Indent {
				primaryRunning = task.DesiredState == "Running" && strings.HasPrefix(task.CurrentState, "Running")
			}
			if task.Indent {
				name = ""
			}
			image := task.Image
			index := strings.Index(task.Image, "@")
			if index != -1 {
				image = task.Image[:index]
			}
			image = sanitizeImageName(image)
			desiredState := task.DesiredState
			if desiredState == "Running" {
				desiredState = ""
			} else if desiredState == "Shutdown" {
				desiredState = ""
			}

			if psAll || (!psAll && primaryRunning && !task.Indent) || !primaryRunning {
				t.AppendRow([]interface{}{
					name,
					task.ID[0:8],
					image,
					task.Node,
					desiredState,
					sanitizeCurrentState(task.CurrentState),
					sanitizeError(task.Error),
				})
			}
			// if task.Error != "" {
			// 	t.AppendRow([]interface{}{
			// 		"",
			// 		"",
			// 		"",
			// 		task.Error,
			// 	})

			// }

			// var yamlData []byte
			// if yamlData, err = yaml.Marshal(task); err != nil {
			// 	return err
			// }
			// fmt.Println(string(yamlData))
		}
		// t.SetAutoIndex(true)
		// t.SetColumnConfigs([]table.ColumnConfig{
		// 	{Number: 4, AutoMerge: true, Align: text.AlignLeft, AlignHeader: text.AlignLeft},
		// 	{Number: 5, AutoMerge: true, Align: text.AlignLeft, AlignHeader: text.AlignLeft},
		// 	{Number: 6, AutoMerge: true, Align: text.AlignLeft, AlignHeader: text.AlignLeft},
		// 	// 	{Number: 4, AutoMerge: true},
		// 	// 	{Number: 5, AutoMerge: true},
		// })
		t.Render()

		return nil
	},
}

func init() {
	serviceCmd.AddCommand(servicePsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	servicePsCmd.Flags().BoolVarP(&psDown, "down", "d", false, "Show only services that are down")
	servicePsCmd.Flags().BoolVarP(&psAll, "all", "a", false, "Show all tasks even if the primary one is running")
}
