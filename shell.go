// Package goshell - A simple shell with piping and redirection
package goshell

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

var shellName = "Go Shell"
var prompt = "G$ "
var startMessage = "Welcome to Go Shell! The shell coded entirely in Go(lang)\n"

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var purple = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

// Run - The run loop. Use for starting the shell
func Run() {
	// This shell is currently Unix only
	if runtime.GOOS == "windows" {
		fmt.Fprintf(os.Stderr, "%s: Sorry - currently unavailable for Windows\n", shellName)
		return
	}

	// Display a start message to the user
	fmt.Println(green + startMessage + reset)

	lineReader := bufio.NewReader(os.Stdin)

	// The main loop of requesting input and processing it
	for {
		fmt.Printf(green+"%s"+reset, prompt)
		cmdLine, err := lineReader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		processLine(cmdLine)
	}
}

// SetPrompt - Change the input prompter to a new string
func SetPrompt(newPrompt string) {
	prompt = newPrompt
}

// SetStartMessage - change the welcome message to a new string
func SetStartMessage(newMessage string) {
	startMessage = newMessage
}

// SetShellName - Set the name of the shell used in errors and more
func SetShellName(newName string) {
	shellName = newName
}

// DisableColors - Disable the use of special colors
func DisableColors() {
	reset = ""
	red = ""
	green = ""
	yellow = ""
	blue = ""
	purple = ""
	cyan = ""
	gray = ""
	white = ""
}
