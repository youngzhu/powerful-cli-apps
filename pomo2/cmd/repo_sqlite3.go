//go:build !inmemory
// +build !inmemory

package cmd

import (
	"github.com/spf13/viper"
	"pomo2/pomodoro"
	"pomo2/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	repo, err := repository.NewSQLite3Repo(viper.GetString("db"))
	if err != nil {
		return nil, err
	}
	return repo, nil
}
