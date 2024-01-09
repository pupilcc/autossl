package middleware

import (
	"autossl/internal/service"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

var usr, _ = user.Current()

func InitAcme() {
	ca()
	email()
	export()
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

func Issue(name string) {
	logger := GetLogger()
	produce := os.Getenv("ACME_PRODUCE")
	alias := os.Getenv("ACME_ALIAS")

	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--issue", "--dns", produce, "-d", name, "--challenge-alias", alias)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Start()
	if err != nil {
		logger.Error("cmd.Start() failed with %s\n", zap.String("error", err.Error()))
	}
}

func Install(name string, id string) {
	logger := GetLogger()
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--install-cert", "-d", name, "--key-file", service.CertPath+id+".key", "--fullchain-file", service.CertPath+id+".crt")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Start()
	if err != nil {
		logger.Error("cmd.Start() failed with %s\n", zap.String("error", err.Error()))
	}
}

func export() {
	logger := GetLogger()
	account := os.Getenv("ACME_ACCOUNT")
	token := os.Getenv("ACME_TOKEN")
	cmd := exec.Command("env")
	cmd.Env = append(cmd.Env, "CF_EMAIL="+account)
	cmd.Env = append(cmd.Env, "CF_KEY="+token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Start()
	if err != nil {
		logger.Error("cmd.Start() failed with %s\n", zap.String("error", err.Error()))
	}
}
