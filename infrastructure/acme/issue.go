package acme

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		cmd = exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--issue", "--dns", dns, "-d", name, "-d", "www."+name, "--challenge-alias", alias, "--keylength", "2048")
	} else {
		cmd = exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--issue", "--dns", dns, "-d", name, "--challenge-alias", alias, "--keylength", "2048")
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
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--install-cert", "-d", name, "--key-file", filepath.Join(CertPath, id+".key"), "--fullchain-file", filepath.Join(CertPath, id+".crt"))
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

func execIssue(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe() running command failed", zap.String("error:", err.Error()))
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("cmd.StderrPipe() running command failed", zap.String("error:", err.Error()))
	}
	if err = cmd.Start(); err != nil {
		logger.Error("cmd.Start() running command failed", zap.String("error:", err.Error()))
	}

	go func() {
		tmp := make([]byte, 1024)
		for {
			n, err := stdout.Read(tmp)
			if n > 0 {
				fmt.Print(string(tmp[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	go func() {
		tmp := make([]byte, 1024)
		for {
			n, err := stderr.Read(tmp)
			if n > 0 {
				fmt.Print(string(tmp[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		logger.Error("cmd.Wait() running command failed", zap.String("error:", err.Error()))
	}
	return err
}
