package pomodoro_test

import (
	"testing"

	"github.com/sashaaKr/golang_command_line/pomo/pomodoro"
	"github.com/sashaaKr/golang_command_line/pomo/pomodoro/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	return repository.NewInMemoryRepo(), func() {}
}
