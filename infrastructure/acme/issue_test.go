package acme

import (
	"os/exec"
	"testing"
)

func TestExecIssueExitStatus(t *testing.T) {
	if err := execIssue(exec.Command("sh", "-c", "exit 1")); err == nil {
		t.Fatal("expected non-cron command failure")
	}

	if err := execIssue(exec.Command("sh", "-c", "exit 1", "--cron")); err != nil {
		t.Fatalf("expected cron exit status 1 to succeed: %v", err)
	}
}
