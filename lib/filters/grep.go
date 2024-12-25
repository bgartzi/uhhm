package filters

import (
	"github.com/bgartzi/uhhm/lib/host"
	"github.com/bgartzi/uhhm/lib/inventory"
	"strings"
)

type GrepHostFilter struct {
	expression string
}

func (grepFilter *GrepHostFilter) GrepName(inventory inventory.InventoryBackend) []host.Host {
	ret := make([]host.Host, 0)
	for _, h := range inventory.ListHosts() {
		if strings.Contains(h.NickName, grepFilter.expression) {
			ret = append(ret, h)
		}
	}
	if len(ret) == 0 {
		return nil
	}
	return ret
}

func (grepFilter *GrepHostFilter) GrepAddress(inventory inventory.InventoryBackend) []host.Host {
	ret := make([]host.Host, 0)
	for _, h := range inventory.ListHosts() {
		if strings.Contains(h.Address, grepFilter.expression) {
			ret = append(ret, h)
		}
	}
	if len(ret) == 0 {
		return nil
	}
	return ret
}

func (grepFilter *GrepHostFilter) GrepInfo(inventory inventory.InventoryBackend) []host.Host {
	ret := make([]host.Host, 0)
	for _, h := range inventory.ListHosts() {
		if strings.Contains(h.Info, grepFilter.expression) {
			ret = append(ret, h)
		}
	}
	if len(ret) == 0 {
		return nil
	}
	return ret
}
