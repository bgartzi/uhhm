package host

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"os/exec"
)

func (host *Host) SimpleClientConfig(knownHostKeyCallback ssh.HostKeyCallback, authMethod ssh.AuthMethod) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: knownHostKeyCallback,
	}
}

func (host *Host) Session(privateKeyFilePath string) error {
	cmd := exec.Command(
		"ssh",
		"-i",
		privateKeyFilePath,
		fmt.Sprintf("ssh://%s@%s", host.User, host.SSHAddress()))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
