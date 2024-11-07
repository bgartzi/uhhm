package local

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bgartzi/uhhmm/lib/host"
	"github.com/bgartzi/uhhmm/lib/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"os"
	"path/filepath"
	"strings"
)

type KnownHostsEntry struct {
	Host    string
	KeyType string
	Key     string
}

func KnownHostsEntryFromKey(host host.Host, key ssh.PublicKey) KnownHostsEntry {
	return KnownHostsEntry{
		Host:    host.Address,
		KeyType: key.Type(),
		Key:     utils.EncodeKeyBase64(key)}
}

func ParseKnownHostsEntry(serialized_entry string) (KnownHostsEntry, error) {
	entry_fields := strings.Fields(serialized_entry)
	if len(entry_fields) != 3 {
		return KnownHostsEntry{}, errors.New("Unknown known_hosts entry format")
	}
	return KnownHostsEntry{
		Host:    entry_fields[0],
		KeyType: entry_fields[1],
		Key:     entry_fields[2]}, nil
}

func (entry KnownHostsEntry) String() string {
	return fmt.Sprintf("%s %s %s\n", entry.Host, entry.KeyType, entry.Key)
}

type KnownHosts struct {
	Path string
}

func KnownHostsFilePath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	knownHostsPath := filepath.Join(userHomeDir, ".ssh", "known_hosts")
	return knownHostsPath, nil
}

func InitKnownHosts() (KnownHosts, error) {
	uhhmKnownHostsPath, err := KnownHostsFilePath()
	if err != nil {
		return KnownHosts{}, err
	}
	if _, err := os.Stat(uhhmKnownHostsPath); errors.Is(err, os.ErrNotExist) {
		// file does not exist
		f_kh, err := os.OpenFile(uhhmKnownHostsPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return KnownHosts{}, err
		}
		defer f_kh.Close()
	}
	return KnownHosts{Path: uhhmKnownHostsPath}, nil
}

func (known_hosts *KnownHosts) AddHost(host host.Host, key ssh.PublicKey) error {
	f, err := os.OpenFile(known_hosts.Path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	entry := KnownHostsEntryFromKey(host, key)
	_, err = f.WriteString(entry.String())
	if err != nil {
		return err
	}
	return nil
}

func (known_hosts *KnownHosts) RemoveHost(host host.Host) error {
	f_original, err := os.Open(known_hosts.Path)
	if err != nil {
		return err
	}
	f_tmp_path := filepath.Join(filepath.Dir(known_hosts.Path), ".tmpKnownHosts")
	f_tmp, err := os.OpenFile(f_tmp_path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	kh_buffer := bufio.NewScanner(f_original)
	tmp_writer := bufio.NewWriter(f_tmp)
	removed := false
	for kh_buffer.Scan() {
		known_host_info := kh_buffer.Text()
		if strings.HasPrefix(known_host_info, host.Address) {
			removed = true
			continue
		}
		fmt.Fprintln(tmp_writer, known_host_info)
	}
	f_original.Close()
	err = tmp_writer.Flush()
	f_tmp.Close()
	if err != nil {
		return err
	}
	if removed {
		os.Remove(known_hosts.Path)
		err = os.Rename(f_tmp_path, known_hosts.Path)
		return err
	}
	err_msg := fmt.Sprintf(
		"Can't remove %s from known_hosts. Not found.", host.Address)
	return errors.New(err_msg)
}

func (knownHost *KnownHosts) KeyCallback() (ssh.HostKeyCallback, error) {
	hostKeyCallback, err := knownhosts.New(knownHost.Path)
	if err != nil {
		return nil, err
	}
	return hostKeyCallback, err
}
