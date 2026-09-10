package acme

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"go.uber.org/zap"
)

var CertPath = "./data/cert"

func Issue(name string) error {
	dns := os.Getenv("ACME_DNS")
	alias := os.Getenv("ACME_ALIAS")

	parts := strings.Split(name, ".")

	var cmd *exec.Cmd
	if len(parts) == 2 {
		cmd = exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--issue", "--dns", dns, "-d", name, "-d", "www."+name, "--challenge-alias", alias, "--keylength", "ec-256")
	} else {
		cmd = exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--issue", "--dns", dns, "-d", name, "--challenge-alias", alias, "--keylength", "ec-256")
	}
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

func Install(name string, id string) error {
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--install-cert", "--ecc", "-d", name, "--key-file", filepath.Join(CertPath, id+".key"), "--fullchain-file", filepath.Join(CertPath, id+".crt"))
	logger.Info("command", zap.String("Running command:", strings.Join(cmd.Args, " ")))
	return execIssue(cmd)
}

func Remove(name string) error {
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--remove", "--ecc", "--domain", name)
	logger.Info("command", zap.String("Running command:", strings.Join(cmd.Args, " ")))
	err := execIssue(cmd)
	if err != nil {
		return err
	}

	folder := filepath.Join(usr.HomeDir, ".acme.sh", name+"_ecc")
	filePaths := []string{folder}
	for _, filePath := range filePaths {
		err := os.RemoveAll(filePath)
		if err != nil {
			return err
		}
	}
	return nil
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
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
