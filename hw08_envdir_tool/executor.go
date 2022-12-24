package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.

func setEnvarmant(env Environment, returnCode *int) {
	var err error = nil
	for key, value := range env {
		if value.NeedRemove {
			err = os.Unsetenv(key)
		} else {
			err = os.Setenv(key, value.Value)
		}
		if err != nil {
			*returnCode = 1
		}
	}
}

func RunCmd(cmd []string, env Environment) (returnCode int) {
	stringCommand := ""

	if len(cmd) >= 2 {
		stringCommand = strings.Join(cmd, " ")
	}

	setEnvarmant(env, &returnCode)

	complete := exec.Command("bash", "-c", stringCommand)
	complete.Stdout = os.Stdout
	complete.Stderr = os.Stderr
	if err := complete.Start(); err != nil {
		returnCode = 1
	}

	if err := complete.Wait(); err != nil {
		returnCode = 1
	}

	return returnCode
}
