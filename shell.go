// The main driver and shell management
package goshell

import (
   "fmt"
   "os"
   "bufio"
   "runtime"
)

var shellName = "Go Shell"
var prompt = "G$ "
var startMessage = "Welcome to Go Shell! The shell coded entirely in Go(lang)\n"

var Reset   = "\033[0m"
var Red     = "\033[31m"
var Green   = "\033[32m"
var Yellow  = "\033[33m"
var Blue    = "\033[34m"
var Purple  = "\033[35m"
var Cyan    = "\033[36m"
var Gray    = "\033[37m"
var White   = "\033[97m"

func Run() {
   if runtime.GOOS == "windows" {
		fmt.Fprintf(os.Stderr, "%s: Sorry - currently unavailable for Windows\n", shellName)
	}
   fmt.Println(startMessage)

   lineReader := bufio.NewReader(os.Stdin)
   for {
      fmt.Printf(Green + "%s" + Reset , prompt)
      cmdLine, err := lineReader.ReadString('\n')
      if err != nil {
         fmt.Fprintln(os.Stderr, err)
      }
      ProcessLine(cmdLine)
   }
}

func SetPrompt(newPrompt string) {
   prompt = newPrompt
}

func SetStartMessage(newMessage string) {
   startMessage = newMessage
}

func SetShellName(newName string) {
   shellName = newName
}

func DisableColors() {
   Reset   = ""
   Red     = ""
   Green   = ""
   Yellow  = ""
   Blue    = ""
   Purple  = ""
   Cyan    = ""
   Gray    = ""
   White   = ""
}