// Parse lines and pipelines
package goshell

import (
	"fmt"
	"errors"
	"strings"
	"os"
)

const (
	inRedirectChar string = "<"
	outRedirectChar string = ">"
	pipeChar string = "|"
)

func ProcessLine(cmdLine string) {
	cmdLine = strings.TrimSpace(cmdLine)
	var commands []command

	if len(cmdLine) == 0 {
		return;
	}

	pipeStages := strings.Split(cmdLine, pipeChar)
	
	commands, err := CreatePipeStages(pipeStages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v.\n", shellName, err)
		return
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

func CreateCommand(pipeStage string) (command, error) {
	args := strings.Fields(pipeStage)
	var cmd command

	executionArgs := make([]string, 0)

	// Deliminating by pipes exposed an empty command
	if len(args) == 0 {
		return cmd, errors.New("Pipeline stage cannot be empty")
	}

	var inRedirectFile, outRedirectFile string

	i := 0
	for i < len(args) {
		if strings.HasPrefix(args[i], inRedirectChar) {
			if i == 0 {
				return cmd, errors.New("Command must preceed input redirection")
			}
			if len(args[i]) > 1 {
				inRedirectFile = args[i][1:]
				i++
			} else if i == (len(args) - 1) {
				return cmd, errors.New("Redirection must include input file name")
			} else {
				inRedirectFile = args[i + 1]
				i += 2
			}
		} else if strings.HasPrefix(args[i], outRedirectChar){
			if i == 0 {
				return cmd, errors.New("Command must preceed output redirection")
			}
			if len(args[i]) > 1 {
				outRedirectFile = args[i][1:]
				i++
			} else if i == (len(args) - 1) {
				return cmd, errors.New("Redirection must include output file name")
			} else {
				outRedirectFile = args[i + 1]
				i += 2
			}
		} else if i < len(args) {
			executionArgs = append(executionArgs, args[i])
			i++
		}
	}
	
	cmd.argsList = executionArgs

	err := cmd.setRedirects(inRedirectFile, outRedirectFile)
	if err != nil {
		return cmd, err
	}

	fmt.Printf("Output Redirection file: %v\n", outRedirectFile)
	fmt.Printf("Input Redirection file: %v\n", inRedirectFile)
	fmt.Printf("Other args: %v\n", executionArgs)

	return cmd, nil
}