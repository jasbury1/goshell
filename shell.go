package goshell

import (
   "fmt"
   "os"
   "bufio"
)

func Display() {
   fmt.Println("This works")
}

func Run() {
   lineReader := bufio.NewReader(os.Stdin)
   for {
      fmt.Print("gosh$ ")
      _, err := lineReader.ReadString('\n')
      if err != nil {
         fmt.Fprintln(os.Stderr, err)
      }
   }
}
