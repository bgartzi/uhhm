package inventory

import (
	"github.com/bgartzi/uhhmm/config"
	"path"
)

func InventoryFilePath() (string, error) {
	uhhmDir, err := config.HomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(uhhmDir, ".HOST_INVENTORY"), nil
}
