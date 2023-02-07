package cmd

import (
	"pomo2/pomodoro"
	"pomo2/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewMemoryRepo(), nil
}
