package commands

import (
	"os/exec"

	"github.com/anircodes/configo/core"
)

func Run(commandName string, commandArguments ...string) (core.Result, error) {

	execOutput := core.ExecOutput{}
	execOutput.ErrorCode = 0
	execOutput.ErrorMessage = ""
	result := core.Result{}
	result.Error = execOutput
	cmd := exec.Command(commandName, commandArguments...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		result.Error.ErrorCode = -1
		result.Error.ErrorMessage = err.Error()
		result.Output = string(out)
		return result, err
	}

	return result, nil

}
