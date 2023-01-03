package main

import "os/exec"

type step struct {
	name    string   // the step name
	exe     string   //  the executable name of the external tool want to execute
	args    []string // the arguments for the executable
	message string   // the output message in case of success
	proj    string   // the target project on which to execute the task
}

func newStep(name, exe, message, proj string, args []string) step {
	return step{
		name:    name,
		exe:     exe,
		message: message,
		proj:    proj,
		args:    args,
	}
}

func (s step) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", newStepErr(s.name, "failed to execute", err)
	}

	return s.message, nil
}
