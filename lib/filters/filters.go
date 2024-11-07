package filters

import (
	"github.com/bgartzi/uhhmm/lib/host"
)

type HostFilter func(inventory []host.Host) []host.Host

type FilterChain struct {
	filters []HostFilter
}

func (chain *FilterChain) Append(filter HostFilter) {
	chain.filters = append(chain.filters, filter)
}

func (chain *FilterChain) Apply(hosts []host.Host) []host.Host {
	for _, filter := range chain.filters {
		hosts = filter(hosts)
	}
	return hosts
}
