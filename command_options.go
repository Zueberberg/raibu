package main

import (
	"fmt"
	"os"
	"strings"
)

type CommandOptions struct {
	ExecName  string
	Args      []string
	HeaderMsg string
	ErrorMsg  string
}

func (o *CommandOptions) ErrorString(e error) string {
	return fmt.Sprintf(o.ErrorMsg, o.ExecName, o.Args, e)
}

func MakeOptions(cmd, headerMsg, errorMsg string) *CommandOptions {
	if cmd != "" {
		cmdargs := strings.Split(cmd, " ")
		if len(cmdargs) < 1 {
			PrintLog("Command is incorrect. Want at least 2 space-separated values.")
			os.Exit(1)
		}

		return &CommandOptions{
			ExecName:  cmdargs[0],
			Args:      cmdargs[1:],
			HeaderMsg: headerMsg,
			ErrorMsg:  errorMsg,
		}
	}
	return nil
}
