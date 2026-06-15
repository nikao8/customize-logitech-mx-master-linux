package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const serviceContent = `[Unit]
Description=Logitech MX Master Configuration Daemon
After=multi-user.target

[Service]
Type=simple
ExecStart=/usr/bin/logid
Restart=always
RestartSec=3
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
`

func GenerateServiceContent() string {
	return serviceContent
}

func WriteConfigFile(content string, path string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func WriteServiceFile(path string) error {
	return os.WriteFile(path, []byte(serviceContent), 0644)
}

func EnableAndStartService() (string, error) {
	var output strings.Builder

	cmd := exec.Command("systemctl", "daemon-reload")
	out, err := cmd.CombinedOutput()
	output.WriteString(fmt.Sprintf("$ systemctl daemon-reload\n%s\n", string(out)))
	if err != nil {
		return output.String(), fmt.Errorf("daemon-reload failed: %w\n%s", err, string(out))
	}

	cmd = exec.Command("systemctl", "enable", "logid.service")
	out, err = cmd.CombinedOutput()
	output.WriteString(fmt.Sprintf("$ systemctl enable logid.service\n%s\n", string(out)))
	if err != nil {
		return output.String(), fmt.Errorf("enable failed: %w\n%s", err, string(out))
	}

	cmd = exec.Command("systemctl", "start", "logid.service")
	out, err = cmd.CombinedOutput()
	output.WriteString(fmt.Sprintf("$ systemctl start logid.service\n%s\n", string(out)))
	if err != nil {
		return output.String(), fmt.Errorf("start failed: %w\n%s", err, string(out))
	}

	return output.String(), nil
}

func StopService() (string, error) {
	cmd := exec.Command("systemctl", "stop", "logid.service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("stop failed: %w\n%s", err, string(out))
	}
	return string(out), nil
}

func RestartService() (string, error) {
	cmd := exec.Command("systemctl", "restart", "logid.service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("restart failed: %w\n%s", err, string(out))
	}
	return string(out), nil
}

func IsServiceRunning() bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "logid.service")
	return cmd.Run() == nil
}

func EnsureConfigDir() error {
	return os.MkdirAll("/etc", 0755)
}
