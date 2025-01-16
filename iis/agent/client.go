package agent

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct {
	Hostname string
	Username string
	Password string
}

func (client Client) Execute(script string) (*[]byte, error) {
	var sb strings.Builder
	sb.WriteString("Invoke-Command ")
	if len(client.Hostname) > 0 {
		sb.WriteString(fmt.Sprintf("-ComputerName '%s' ", client.Hostname))
	}

	if len(client.Username) > 0 && len(client.Password) > 0 {
		sb.WriteString(fmt.Sprintf(`-Credential (New-Object System.Management.Automation.PSCredential ('%s', (ConvertTo-SecureString '%s' -AsPlainText -Force))) -Authentication Negotiate `, client.Username, client.Password))
	}

	script = strings.ReplaceAll(script, `"`, `'`)
	sb.WriteString(fmt.Sprintf("-ScriptBlock { param() %s } ", script))
	command := append([]string{"-NoProfile", "-NonInteractive"}, sb.String())

	ps, _ := exec.LookPath("powershell.exe")
	cmd := exec.Command(ps, command...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		if stderr.Len() > 0 {
			fmt.Printf("Stderr: %s\n", stderr.String())
			err = errors.New(string(stderr.Bytes()))
		}

		fmt.Printf("Stdout: %s\n", stdout.String())
		return nil, err
	}

	if stderr.Len() > 0 {
		err = errors.New(string(stderr.Bytes()))
	}

	bytes := stdout.Bytes()
	return &bytes, nil
}
