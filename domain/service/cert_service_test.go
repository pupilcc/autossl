package service

import "testing"

func TestBaseDomain(t *testing.T) {
	for input, want := range map[string]string{
		"latticeway.com":      "latticeway.com",
		"*.latticeway.com":    "latticeway.com",
		" *.LATTICEWAY.COM. ": "latticeway.com",
	} {
		if got := baseDomain(input); got != want {
			t.Errorf("baseDomain(%q) = %q, want %q", input, got, want)
		}
	}
}
