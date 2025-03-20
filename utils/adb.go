package utils

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
)

type AdbUtils struct {
}

func StartAdbDeamon() {
	cmd := exec.Command("adb", "start-server")

	err := cmd.Run()
	if err != nil {
		log.Panic("failed to start adb server: %w", err)
	}
}
