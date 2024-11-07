package utils

import (
	"encoding/base64"
	"golang.org/x/crypto/ssh"
)

func EncodeKeyBase64(key ssh.PublicKey) string {
	return base64.StdEncoding.EncodeToString(key.Marshal())
}
