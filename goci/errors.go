package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation failed")
)

type stepErr struct {
	step  string // record the step name
	msg   string // describes the condition
	cause error  // the underlying error
}

func (s *stepErr) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.step,
		s.msg, s.cause)
}

func (s *stepErr) Is(target error) bool {
	t, ok := target.(*stepErr)
	if !ok {
		return false
	}
	return t.step == s.step
}

func (s *stepErr) Unwrap() error {
	return s.cause
}

func newStepErr(step, msg string, cause error) *stepErr {
	return &stepErr{
		step:  step,
		msg:   msg,
		cause: cause,
	}
}
