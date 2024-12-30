package local

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/bgartzi/uhhm/config"
	"github.com/bgartzi/uhhm/host"
	"github.com/bgartzi/uhhm/utils"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func UHHMPrivateSSHKeyPath() (string, error) {
	uhhmHomeDir, err := config.HomeDir()
	if err != nil {
		return "", err
	}
	// uhhmKnownHostsPath := filepath.Join(uhhmHomeDir, ".UHHM_KNOWN_HOSTS")
	return filepath.Join(uhhmHomeDir, ".UHHM_RSA_KEY"), nil
}

func PublicKeyPath(privateKeyPath string) string {
	return privateKeyPath + ".pub"
}

func CreateRSAKeyPair(privateKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}

	privateKeyFile, err := os.OpenFile(privateKeyPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(
		PublicKeyPath(privateKeyPath),
		ssh.MarshalAuthorizedKey(publicKey),
		0644)
}

// Create a new RSA key pair if it doesn't already exist
func SafeCreateRSAKeyPair(privateKeyPath string) error {
	if utils.FileExists(privateKeyPath) {
		return nil
	}
	return CreateRSAKeyPair(privateKeyPath)
}

func PrivateKeyAuthMeth(privateKeyPath string) (ssh.AuthMethod, error) {
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(privateKey), nil
}

func CopyPublicKeyTo(targetHost host.Host, privateKeyPath string) error {
	fmt.Printf("Copying UHHM rsa public key to target host. It might ask for remote host's user password.\n")
	cmd := exec.Command(
		"ssh-copy-id",
		"-i",
		PublicKeyPath(privateKeyPath),
		fmt.Sprintf("%s@%s", targetHost.User, targetHost.Address))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
