package acme

import (
	"autossl/config"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

var usr, _ = user.Current()
var logger = config.GetLogger()

func InitAcme() {
	ca()
	email()
	upgrade()
}

func ca() {
	ca := os.Getenv("ACME_CA")
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--set-default-ca", "--server", ca)
	execSetting(cmd)
}

func email() {
	email := os.Getenv("ACME_EMAIL")
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--update-account", "--email", email)
	execSetting(cmd)
}

func upgrade() {
	cmd := exec.Command(filepath.Join(usr.HomeDir, ".acme.sh/acme.sh"), "--upgrade", "--auto-upgrade")
	execSetting(cmd)
}

func execSetting(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Start()
	if err != nil {
		logger.Error("cmd.Start() failed: ", zap.String("error", err.Error()))
	}
}
