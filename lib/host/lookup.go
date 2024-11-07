package host

import (
	"slices"
)

func (host *Host) HasLabel(label string) bool {
	return slices.Contains(host.Labels, label)
}

func (host *Host) HasLabels(labels []string) bool {
	for _, label := range labels {
		if !host.HasLabel(label) {
			return false
		}
	}
	return true
}

func (host *Host) HasAnyLabel(labels []string) bool {
	for _, label := range labels {
		if host.HasLabel(label) {
			return true
		}
	}
	return false
}
