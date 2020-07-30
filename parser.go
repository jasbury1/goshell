// Package goshell - A simple shell with piping and redirection
package goshell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	inRedirectChar  string = "<"
	outRedirectChar string = ">"
	pipeChar        string = "|"
)

// Process a line of input and execute it
func processLine(cmdLine string) {
	cmdLine = strings.TrimSpace(cmdLine)

	commandList := make([]exec.Cmd, 0)

	if len(cmdLine) == 0 {
		return
	}

	pipeStages := strings.Split(cmdLine, pipeChar)

	err := createPipeStages(&commandList, pipeStages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v.\n", shellName, err)
		return
	}

	err = connectPipeline(commandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: Error with pipes: %v.\n", shellName, err)
		return
	}

	err = executePipeline(commandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: Error during execution: %v\n", shellName, err)
		return
	}
}

// Separate out a string slice of pipeline of commands into Cmd slice
func createPipeStages(commandList *[]exec.Cmd, stages []string) error {

	for _, stage := range stages {
		cmd, err := createCommand(stage)
		if err != nil {
			return err
		}
		*commandList = append(*commandList, *cmd)
	}
	return nil
}

// Creates a command using a string of one pipeline step
func createCommand(pipeStage string) (*exec.Cmd, error) {
	args := strings.Fields(pipeStage)
	executionArgs := make([]string, 0)
	var cmd *exec.Cmd

	// Deliminating by pipes exposed an empty command
	if len(args) == 0 {
		return nil, errors.New("Pipeline stage cannot be empty")
	}

	var inRedirectFile, outRedirectFile string

	// Any redirection specifiers (</>) will get parsed and saved
	i := 0
	for i < len(args) {
		if strings.HasPrefix(args[i], inRedirectChar) {
			if i == 0 {
				return nil, errors.New("Command must precede input redirection")
			}
			if len(args[i]) > 1 {
				inRedirectFile = args[i][1:]
				i++
			} else if i == (len(args) - 1) {
				return nil, errors.New("Redirection must include input file name")
			} else {
				inRedirectFile = args[i+1]
				i += 2
			}
		} else if strings.HasPrefix(args[i], outRedirectChar) {
			if i == 0 {
				return nil, errors.New("Command must precede output redirection")
			}
			if len(args[i]) > 1 {
				outRedirectFile = args[i][1:]
				i++
			} else if i == (len(args) - 1) {
				return nil, errors.New("Redirection must include output file name")
			} else {
				outRedirectFile = args[i+1]
				i += 2
			}
		} else if i < len(args) {
			// Save command arguments only if they arent references to redirection
			executionArgs = append(executionArgs, args[i])
			i++
		}
	}
	// Create a command using only the command name and arguments
	cmd = exec.Command(executionArgs[0], executionArgs[1:]...)

	// Set the redirect output files from the data we parsed
	err := setRedirects(cmd, inRedirectFile, outRedirectFile)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// Takes in strings of any files for redirection. Opens them and associates
// them with the Cmd object as a File* (Fd)
func setRedirects(cmd *exec.Cmd, inFile string, outFile string) error {
	if inFile != "" {
		inF, err := os.OpenFile(inFile, syscall.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		cmd.Stdin = inF
	}
	if outFile != "" {
		outF, err := os.OpenFile(outFile, (syscall.O_WRONLY | syscall.O_CREAT), 0666)
		if err != nil {
			return err
		}
		cmd.Stdout = outF
	}
	return nil
}

// Connect the commands together as a variable length pipeline
func connectPipeline(commandList []exec.Cmd) error {
	var err error

	for i := range commandList {
		if i == len(commandList)-1 {
			break
		}
		if commandList[i+1].Stdin != nil || commandList[i].Stdout != nil {
			return errors.New("Ambiguous input for file redirection and pipe")
		}
		commandList[i+1].Stdin, err = commandList[i].StdoutPipe()
		if err != nil {
			return err
		}
	}
	if commandList[0].Stdin == nil {
		commandList[0].Stdin = os.Stdin
	}
	if commandList[len(commandList)-1].Stdout == nil {
		commandList[len(commandList)-1].Stdout = os.Stdout
	}
	return nil
}

// Execute all commands. Pipeline commands will wait on the previous command
func executePipeline(commandList []exec.Cmd) error {
	//TODO: Handle case where the command is built in
	// Start execution in reverse pipe order
	for i := len(commandList) - 1; i > 0; i-- {
		commandList[i].Start()
	}

	// Make the commands wait for completion
	for i := range commandList {
		if i == 0 {
			err := commandList[i].Run()
			if err != nil {
				return err
			}
		} else {
			err := commandList[i].Wait()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
