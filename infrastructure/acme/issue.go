package acme

import (
	"autossl/domain/service"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Issue(name string) error {
	dns := os.Getenv("ACME_DNS")
	alias := os.Getenv("ACME_ALIAS")

	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--issue", "--dns", dns, "-d", name, "--challenge-alias", alias, "--keylength", "2048")
	return execIssue(cmd)
}

func Install(name string, id string) error {
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--install-cert", "-d", name, "--key-file", filepath.Join(service.CertPath, id+".key"), "--fullchain-file", filepath.Join(service.CertPath, id+".crt"))
	logger.Info("command", zap.String("Running command:", strings.Join(cmd.Args, " ")))
	return execIssue(cmd)
}

func Remove(name string) error {
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--remove", "--domain", name)
	logger.Info("command", zap.String("Running command:", strings.Join(cmd.Args, " ")))
	err := execIssue(cmd)
	if err != nil {
		return err
	}

	folder := filepath.Join("./.acme.sh", name)
	filePaths := []string{folder}
	for _, filePath := range filePaths {
		err := os.Remove(filePath)
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

func execIssue(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe() running command failed", zap.String("error:", err.Error()))
	}
	if err = cmd.Start(); err != nil {
		logger.Error("cmd.Start() running command failed", zap.String("error:", err.Error()))
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
	if err := cmd.Wait(); err != nil {
		logger.Error("cmd.Wait() running command failed", zap.String("error:", err.Error()))
	}
	return err
}
