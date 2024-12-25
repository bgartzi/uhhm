package inventory

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/bgartzi/uhhm/lib/host"
	"io/fs"
	"os"
)

const defaultInventoryPerms fs.FileMode = 0600

type simple_inventory struct {
	FilePath string
	Version  string
	Hosts    []host.Host
}

func get_empty_inventory(file_path string) *simple_inventory {
	return &simple_inventory{FilePath: file_path, Version: "0.0"}
}

func (inventory *simple_inventory) Read() ([]host.Host, error) {
	version := inventory.Version
	f, err := os.OpenFile(inventory.FilePath, os.O_RDONLY, defaultInventoryPerms)
	if err != nil {
		f.Close()
		err = inventory.Write()
		return inventory.Hosts, err
	}
	dec := gob.NewDecoder(f)
	err = dec.Decode(inventory)
	if err != nil {
		// FIXME: Inventory file is empty or corrupted, write empty sorry
		f.Close()
		err = inventory.Write()
		return inventory.Hosts, err
	}
	if version != inventory.Version {
		fmt.Printf(
			"WARNING: Inventory file (%s) and implementation (%s) versions mismatch.\n",
			inventory.Version,
			version)
		inventory.Version = version
	}
	f.Close()
	return inventory.Hosts, nil
}

func (inventory *simple_inventory) Write() error {
	f, err := os.OpenFile(inventory.FilePath, os.O_WRONLY|os.O_CREATE, defaultInventoryPerms)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(f)
	err = enc.Encode(inventory)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func (inventory *simple_inventory) find(host host.Host) int {
	for i, h := range inventory.Hosts {
		if h.Address == host.Address || h.NickName == host.NickName {
			return i
		}
	}
	return -1
}

func (inventory *simple_inventory) Contains(host host.Host) bool {
	return inventory.find(host) != -1
}

func (inventory *simple_inventory) AddHost(host host.Host) (int, error) {
	id := inventory.find(host)
	if id != -1 {
		h := inventory.Hosts[id]
		error_msg := fmt.Sprintf(
			"Host %d, %s (%s), is already in the inventory.\n",
			id, h.NickName, h.Address)
		return id, errors.New(error_msg)
	}
	inventory.Hosts = append(inventory.Hosts, host)
	return len(inventory.Hosts) - 1, nil
}

func (inventory *simple_inventory) GetHost(id int) (host.Host, error) {
	if id >= len(inventory.Hosts) {
		err_msg := fmt.Sprintf("Invalid Host ID %d. Out of range.\n", id)
		return host.Host{}, errors.New(err_msg)
	}
	return inventory.Hosts[id], nil
}

func (inventory *simple_inventory) DelHost(id int) (host.Host, error) {
	removed_host, err := inventory.GetHost(id)
	if err != nil {
		return removed_host, err
	}
	new_hosts := make([]host.Host, 0)
	new_hosts = append(new_hosts, inventory.Hosts[:id]...)
	new_hosts = append(new_hosts, inventory.Hosts[id+1:]...)
	inventory.Hosts = new_hosts
	return removed_host, nil
}

func (inventory *simple_inventory) SearchHost(nickname string) (host.Host, error) {
	for _, h := range inventory.Hosts {
		if h.NickName == nickname {
			return h, nil
		}
	}
	err_msg := fmt.Sprintf(
		"Nickname '%s' is not registered in the inventory.\n",
		nickname)
	return host.Host{}, errors.New(err_msg)
}

func (inventory *simple_inventory) ListHosts() []host.Host {
	return inventory.Hosts
}

func InitSimpleInventory(file_path string) *simple_inventory {
	//FIXME: hardcoded semver lol
	inventory := get_empty_inventory(file_path)
	inventory.Read()
	return inventory
}
