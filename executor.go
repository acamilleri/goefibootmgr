package goefibootmgr

import (
	"os/exec"
)

// ExecCommand represent executor who execute the given command
type ExecCommand interface {
	Run(command string, args ...string) error
	Output(command string, args ...string) ([]byte, error)
}

type defaultExecutor struct{}

// Run execute command
func (defaultExecutor) Run(command string, args ...string) error {
	return exec.Command(command, args...).Run()
}

// Output execute command and returned output of command
func (defaultExecutor) Output(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).Output()
}
