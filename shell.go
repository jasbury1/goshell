// The main driver and shell management
package goshell

import (
   "fmt"
   "os"
   "bufio"
)

var shellName = "Go Shell"
var prompt = "G$ "
var startMessage = "Welcome to Go Shell! The shell coded entirely in Go(lang)\n"

func Run() {
   fmt.Println(startMessage)

   lineReader := bufio.NewReader(os.Stdin)
   for {
      fmt.Printf("%s", prompt)
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
