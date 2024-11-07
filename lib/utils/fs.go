package utils

import (
	"os"
)

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	// If it is not a dir, it is a file :s
	return !info.IsDir()
}
