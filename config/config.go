package config

import (
	"errors"
	"fmt"
	"os"
	"path"
)

func ensureDirExists(path string) error {
	pathInfo, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			return err
		}
	} else if !pathInfo.IsDir() {
		return fmt.Errorf("%s exists and is not a directory", path)
	}
	return nil
}

func HomeDir() (string, error) {
	uhhmDir, found := os.LookupEnv("UHHM_HOME")
	if !found || uhhmDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("Can't retrieve user home dir.")
		}
		uhhmDir = path.Join(homeDir, ".uhhm")
	}
	err := ensureDirExists(uhhmDir)
	if err != nil {
		return "", err
	}
	return uhhmDir, nil
}
