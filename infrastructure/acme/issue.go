package acme

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"go.uber.org/zap"
)

var CertPath = "./data/cert"

var ErrIncorrectTXTRecord = errors.New("DNS verification returned an incorrect TXT record")

func Issue(name string) error {
	dns := os.Getenv("ACME_DNS")
	alias := os.Getenv("ACME_ALIAS")
	if strings.EqualFold(os.Getenv("ACME_ALIAS_PER_DOMAIN"), "true") {
		alias = perDomainAlias(name, alias)
	}

	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), issueArgs(name, dns, alias)...)
	if strings.EqualFold(os.Getenv("ACME_DEBUG"), "true") {
		cmd.Args = append(cmd.Args, "--debug", "1")
	}

	logger.Info("command", zap.String("Running command:", strings.Join(cmd.Args, " ")))
	err := execIssue(cmd)
	if err != nil {
		logger.Error("acme.Issue() running command failed", zap.String("error:", err.Error()))
		return err
	}
	return nil
}

func issueArgs(name string, dns string, alias string) []string {
	return []string{"--issue", "--dns", dns, "-d", name, "-d", "*." + name, "--challenge-alias", alias, "--challenge-alias", alias, "--keylength", "ec-256"}
}

func perDomainAlias(name string, base string) string {
	domain := strings.TrimSuffix(strings.ToLower(name), ".")
	base = strings.TrimSuffix(strings.ToLower(base), ".")
	return strings.TrimPrefix(domain, "*.") + "." + base
}

func Install(name string, id string) error {
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--install-cert", "--ecc", "-d", name, "--key-file", filepath.Join(CertPath, id+".key"), "--fullchain-file", filepath.Join(CertPath, id+".crt"))
	logger.Info("command", zap.String("Running command:", strings.Join(cmd.Args, " ")))
	return execIssue(cmd)
}

func Remove(name string) error {
	return os.RemoveAll(filepath.Join(usr.HomeDir, ".acme.sh", name+"_ecc"))
}

func Cron() error {
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--cron", "--home", filepath.Join(usr.HomeDir, ".acme.sh"))
	logger.Info("command", zap.String("Running command:", strings.Join(cmd.Args, " ")))
	return execIssue(cmd)
}

// isCronOperation checks if the command is a cron operation
func isCronOperation(cmd *exec.Cmd) bool {
	return slices.Contains(cmd.Args, "--cron")
}

func execIssue(cmd *exec.Cmd) error {
	header, err := os.CreateTemp("", "autossl-acme-header-*")
	if err != nil {
		return err
	}
	_ = header.Close()
	defer os.Remove(header.Name())

	cmd.Env = append(cmd.Environ(), "HTTP_HEADER="+header.Name())
	var output bytes.Buffer
	stream := io.MultiWriter(os.Stdout, &output)
	cmd.Stdout = stream
	cmd.Stderr = stream
	err = cmd.Run()
	if err != nil {
		if strings.Contains(output.String(), "Incorrect TXT record") {
			return ErrIncorrectTXTRecord
		}
		// Check if this is a cron operation and exit status is 1
		// acme.sh returns exit status 1 when no certificates need renewal during cron
		if isCronOperation(cmd) {
			if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
				logger.Info("acme cron completed - no certificates needed renewal")
				return nil
			}
		}
		logger.Error("cmd.Run() running command failed", zap.String("error:", err.Error()))
	}
	return err
}
