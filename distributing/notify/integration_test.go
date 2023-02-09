//go:build integration
// +build integration

package notify_test

import (
	"distributing/notify"
	"testing"
)

func TestNotify_Send(t *testing.T) {
	n := notify.New("test title", "test msg", notify.SeverityNormal)

	err := n.Send()
	if err != nil {
		t.Error(err)
	}
}
