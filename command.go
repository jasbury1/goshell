package goshell

import (
	"syscall"
	"os"
)

type command struct {
	argsList []string
	inFd int
	outFd int
}

func (c *command) setRedirects(inFile string, outFile string) error {
	curPath, err := os.Getwd()
	if err != nil {
		return err
	}
	if inFile != ""{
		inFd, err := syscall.Open(curPath + "/" + inFile, syscall.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		c.inFd = inFd
	}
	if outFile != "" {
		outFd, err := syscall.Open(curPath + "/" + outFile, (syscall.O_WRONLY | syscall.O_CREAT), 0666)
		if err != nil {
			return err
		}
		c.outFd = outFd
	}
	return nil
}