package host

import (
	"golang.org/x/crypto/ssh"
	"net"
)

func (host Host) PublicKey() (ssh.PublicKey, error) {
	var host_key ssh.PublicKey

	config := ssh.ClientConfig{
		User: host.User,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			host_key = key
			return nil
		},
	}
	conn, err := ssh.Dial("tcp", host.SSHAddress(), &config)
	if err == nil {
		defer conn.Close()
	}
	if host_key == nil {
		return nil, err
	}
	return host_key, nil
}
