// Parse lines and pipelines
package goshell

import (
	"fmt"
	"errors"
	"strings"
	"os"
	"syscall"
	"os/exec"
)

const (
	inRedirectChar string = "<"
	outRedirectChar string = ">"
	pipeChar string = "|"
)

func ProcessLine(cmdLine string) {
	cmdLine = strings.TrimSpace(cmdLine)
	
	commandList := make([]exec.Cmd, 0)

	if len(cmdLine) == 0 {
		return;
	}

	pipeStages := strings.Split(cmdLine, pipeChar)
	
	err := CreatePipeStages(&commandList, pipeStages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v.\n", shellName, err) 
	}

	fmt.Printf("%v\n", commandList)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v.\n", shellName, err)
		return
	}
	
	err = ConnectPipeline(commandList)
	fmt.Printf("%v\n", commandList)	

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: Error with pipes: %v.\n", shellName, err)
		return
	}

	err = ExecutePipeline(commandList)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: Error during execution: %v\n", shellName, err)
		return
	}
}

func CreatePipeStages(commandList *[]exec.Cmd, stages []string) error {
	
	for _, stage := range stages {
		cmd, err := CreateCommand(stage)
		if err != nil {
			return err
		}
		*commandList = append(*commandList, *cmd)

	}
	/*
	if len(pipeline) != 1 {
		err := checkvalidRedirects(pipeline)
		if err != nil {
			return nil, err
		}
	}
	*/

	return nil
}

func CreateCommand(pipeStage string) (*exec.Cmd, error) {
	args := strings.Fields(pipeStage)
	executionArgs := make([]string, 0)
	var cmd *exec.Cmd

	// Deliminating by pipes exposed an empty command
	if len(args) == 0 {
		return nil, errors.New("Pipeline stage cannot be empty")
	}

	var inRedirectFile, outRedirectFile string

	i := 0
	for i < len(args) {
		if strings.HasPrefix(args[i], inRedirectChar) {
			if i == 0 {
				return nil, errors.New("Command must preceed input redirection")
			}
			if len(args[i]) > 1 {
				inRedirectFile = args[i][1:]
				i++
			} else if i == (len(args) - 1) {
				return nil, errors.New("Redirection must include input file name")
			} else {
				inRedirectFile = args[i + 1]
				i += 2
			}
		} else if strings.HasPrefix(args[i], outRedirectChar){
			if i == 0 {
				return nil, errors.New("Command must preceed output redirection")
			}
			if len(args[i]) > 1 {
				outRedirectFile = args[i][1:]
				i++
			} else if i == (len(args) - 1) {
				return nil, errors.New("Redirection must include output file name")
			} else {
				outRedirectFile = args[i + 1]
				i += 2
			}
		} else if i < len(args) {
			executionArgs = append(executionArgs, args[i])
			i++
		}
	}
	cmd = exec.Command(executionArgs[0], executionArgs[1:]...)
	
	err := setRedirects(cmd, inRedirectFile, outRedirectFile)
	
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func setRedirects(cmd *exec.Cmd, inFile string, outFile string) error {
	if inFile != ""{
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

/*
func checkvalidRedirects(pipeline []command) error {
	for i, cmd := range pipeline {
		if i != 0 && cmd.inFd != 0 {
			return errors.New("Ambiguous input for file redirection and pipe")
		}
		if i != len(pipeline) - 1 && cmd.outFd != 0 {
			return errors.New("Ambiguous output for file redirection and pipe")
		}
	}
	return nil
}
*/

func ConnectPipeline(commandList []exec.Cmd) error {
	var err error

	for i, _ := range commandList {
		if i == len(commandList) - 1 {
			break
		}
		fmt.Println("Pipe")

		commandList[i + 1].Stdin, err = commandList[i].StdoutPipe()
		if err != nil {
			return err
		}
	}
	if commandList[0].Stdin == nil {
		commandList[0].Stdin = os.Stdin
	}
	if commandList[len(commandList) - 1].Stdout == nil {
		commandList[len(commandList) - 1].Stdout = os.Stdout
	}
	return nil
}

func ExecutePipeline(commandList []exec.Cmd) error {
	//TODO: Handle case where the command is built in
	// Start execution in reverse pipe order
	for i := len(commandList) - 1; i > 0; i-- {
		commandList[i].Start()
	}

	// Make the commands wait for completion
	for i,_ := range commandList {
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