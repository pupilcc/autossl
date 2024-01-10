package middleware

import (
	"autossl/internal/service"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

var usr, _ = user.Current()

func InitAcme() {
	ca()
	email()
	upgrade()
}

func ca() {
	logger := GetLogger()
	ca := os.Getenv("ACME_CA")
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--set-default-ca", "--server", ca)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Start()
	if err != nil {
		logger.Error("cmd.Start() failed with %s\n", zap.String("error", err.Error()))
	}
}

func email() {
	logger := GetLogger()
	email := os.Getenv("ACME_EMAIL")
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--update-account", "--email", email)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Start()
	if err != nil {
		logger.Error("cmd.Start() failed with %s\n", zap.String("error", err.Error()))
	}
}

func upgrade() {
	logger := GetLogger()
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--upgrade", "--auto-upgrade")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Start()
	if err != nil {
		logger.Error("cmd.Start() failed with %s\n", zap.String("error", err.Error()))
	}
}

func Issue(name string) error {
	logger := GetLogger()
	dns := os.Getenv("ACME_DNS")
	alias := os.Getenv("ACME_ALIAS")

	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--issue", "--dns", dns, "-d", name, "--challenge-alias", alias, "--keylength", "2048")
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

func Install(name string, id string) error {
	logger := GetLogger()

	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--install-cert", "-d", name, "--key-file", filepath.Join(service.CertPath, id+".key"), "--fullchain-file", filepath.Join(service.CertPath, id+".crt"))
	logger.Info("command", zap.String("Running command:", strings.Join(cmd.Args, " ")))
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
