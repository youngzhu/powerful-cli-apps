package main

import (
	"context"
	"os/exec"
	"time"
)

type timeoutStep struct {
	step
	timeout time.Duration
}

func newTimeoutStep(name, exe, message, proj string,
	args []string, timeout time.Duration) timeoutStep {
	s := timeoutStep{}
	s.step = newStep(name, exe, message, proj, args)
	s.timeout = timeout
	if s.timeout == 0 {
		s.timeout = 30 * time.Second
	}
	return s
}

func (t timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, t.exe, t.args...)
	cmd.Dir = t.proj
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &stepErr{
				step:  t.name,
				msg:   "failed time out",
				cause: context.DeadlineExceeded,
			}
		}
		return "", &stepErr{
			step:  t.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return t.message, nil
}
