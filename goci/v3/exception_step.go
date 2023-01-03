package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type exceptionStep struct {
	step
}

func newExceptionStep(name, exe, message, proj string,
	args []string) exceptionStep {
	s := exceptionStep{}
	s.step = newStep(name, exe, message, proj, args)
	return s
}

func (e exceptionStep) execute() (string, error) {
	cmd := exec.Command(e.exe, e.args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = e.proj
	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  e.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	if out.Len() > 0 {
		return "", &stepErr{
			step:  e.name,
			msg:   fmt.Sprintf("invalid format: %s", out.String()),
			cause: nil,
		}
	}

	return e.message, nil
}
