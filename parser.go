// Parse lines and pipelines
package goshell

import (
	"fmt"
	"errors"
	"strings"
	"os"
)

func ProcessLine(cmdLine string) {
	fmt.Println("test")
	pipeStages := strings.Split(cmdLine, "|")
	err := ValidatePipeStages(pipeStages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", shellName, err)
	}
}

func ValidatePipeStages(stages []string) error{
	return errors.New("Invalid pipeline")
}