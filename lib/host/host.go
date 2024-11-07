package host

import (
	"fmt"
)

type Host struct {
	Address  string
	Port     string
	NickName string
	Info     string
	User     string
	Labels   []string
}

func (host *Host) SSHAddress() string {
	return fmt.Sprintf("%s:%s", host.Address, host.Port)
}
