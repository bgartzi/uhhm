package inventory

import (
	"github.com/bgartzi/uhhmm/host"
)

type InventoryBackend interface {
	Read() ([]host.Host, error)
	Write() error
	Contains(host host.Host) bool
	AddHost(host host.Host) (int, error)
	DelHost(id int) (host.Host, error)
	Listhosts() []host.Host
	GetHost(id int) (host.Host, error)
	SearchHost(nickname string) (host.Host, error)
	ListHosts() []host.Host
}
