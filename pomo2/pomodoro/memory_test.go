//go:build inmemory
// +build inmemory

package pomodoro_test

import (
	"pomo2/pomodoro"
	"pomo2/pomodoro/repository"
	"testing"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	return repository.NewMemoryRepo(), func() {}
}
