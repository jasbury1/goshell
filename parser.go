// Parse lines and pipelines
package goshell

import (
	"fmt"
	"errors"
	"strings"
	"os"
)

type command struct {
	outRedirects []string
	inRedirects []string
}

func ProcessLine(cmdLine string) {
	cmdLine = strings.TrimSpace(cmdLine)
	var commands []command

	if len(cmdLine) == 0 {
		return;
	}

	pipeStages := strings.Split(cmdLine, "|")
	
	commands, err := CreatePipeStages(pipeStages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v.\n", shellName, err)
	}
	
	fmt.Printf("%v\n", commands)
}

func CreatePipeStages(stages []string) ([]command, error) {
	pipeline := make([]command, 0)
	for _, stage := range stages {
		cmd, err := CreateCommand(stage)
		if err != nil {
			return nil, err
		}
		pipeline = append(pipeline, cmd)
	}
	return pipeline, nil
}

func CreateCommand(cmd string) (command, error) {
	args := strings.Fields(cmd)
	// Deliminating by pipes exposed an empty command
	if len(args) == 0 {
		return command{}, errors.New("Pipeline stage cannot be empty")
	}
	return command{"test", "test"}, nil
}