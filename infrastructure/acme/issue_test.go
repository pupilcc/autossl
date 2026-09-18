package acme

import (
	"errors"
	"os/exec"
	"slices"
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

func TestExecIssueIncorrectTXTRecord(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo 'Verification error details: Incorrect TXT record'; exit 1")
	if err := execIssue(cmd); !errors.Is(err, ErrIncorrectTXTRecord) {
		t.Fatalf("execIssue() error = %v, want %v", err, ErrIncorrectTXTRecord)
	}
}

func TestPerDomainAlias(t *testing.T) {
	tests := map[string]string{
		"*.latticeway.com": "latticeway.com.autossl.in",
		"latticeway.com":   "latticeway.com.autossl.in",
		"*.Example.COM.":   "example.com.autossl.in",
	}

	for domain, want := range tests {
		if got := perDomainAlias(domain, "AUTOSSL.IN."); got != want {
			t.Errorf("perDomainAlias(%q) = %q, want %q", domain, got, want)
		}
	}
}

func TestIssueArgs(t *testing.T) {
	got := issueArgs("latticeway.com", "dns_cf", "latticeway.com.autossl.in")
	want := []string{"--issue", "--dns", "dns_cf", "-d", "latticeway.com", "-d", "*.latticeway.com", "--challenge-alias", "latticeway.com.autossl.in", "--challenge-alias", "latticeway.com.autossl.in", "--keylength", "ec-256"}
	if !slices.Equal(got, want) {
		t.Fatalf("issueArgs() = %q, want %q", got, want)
	}
}

func TestExecIssueSetsHTTPHeader(t *testing.T) {
	cmd := exec.Command("sh", "-c", "test -f \"$HTTP_HEADER\" && test \"$AUTOSSL_TEST_ENV\" = kept")
	cmd.Env = append(cmd.Environ(), "AUTOSSL_TEST_ENV=kept")
	if err := execIssue(cmd); err != nil {
		t.Fatalf("expected HTTP_HEADER to reference a temporary file: %v", err)
	}
}
