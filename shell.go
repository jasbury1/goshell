// The main driver and shell management. See LICENSE
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

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// The run loop. Use for starting the shell
func Run() {
	// This shell is currently Unix only
	if runtime.GOOS == "windows" {
		fmt.Fprintf(os.Stderr, "%s: Sorry - currently unavailable for Windows\n", shellName)
		return
	}

	// Display a start message to the user
	fmt.Println(Green + startMessage + Reset)

	lineReader := bufio.NewReader(os.Stdin)

	// The main loop of requesting input and processing it
	for {
		fmt.Printf(Green+"%s"+Reset, prompt)
		cmdLine, err := lineReader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		processLine(cmdLine)
	}
}

// Change the input prompter to a new string
func SetPrompt(newPrompt string) {
	prompt = newPrompt
}

// Change the welcome message to a new string
func SetStartMessage(newMessage string) {
	startMessage = newMessage
}

// Set the name of the shell. Used in debugging
func SetShellName(newName string) {
	shellName = newName
}

// Disable the use of special colors
func DisableColors() {
	Reset = ""
	Red = ""
	Green = ""
	Yellow = ""
	Blue = ""
	Purple = ""
	Cyan = ""
	Gray = ""
	White = ""
}
