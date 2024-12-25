package inventory

import (
	"github.com/bgartzi/uhhm/lib/uhhm"
	"path"
)

func InventoryFilePath() (string, error) {
	uhhmDir, err := uhhm.HomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(uhhmDir, ".HOST_INVENTORY"), nil
}
