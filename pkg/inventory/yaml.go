package inventory

import (
	"errors"
	"fmt"
	"io/ioutil"
	"pneuma/pkg/device"

	"gopkg.in/yaml.v3"
)

type YamlInventory struct {
	Filename string
}

func (i *YamlInventory) NewInventory() (inv *Inventory, err error) {
	inv = &Inventory{
		Devices: make(map[string]*device.Device),
		Groups:  make(map[string]*Group),
	}
	b, err := ioutil.ReadFile(i.Filename)
	if err != nil {
		return inv, errors.New("unable to read file")
	}
	err = yaml.Unmarshal(b, inv)
	if err != nil {
		return inv, errors.New("unable to unmarshal YAML")
	}
	fmt.Println("merging")
	_ = inv.MergeGroupsAndDevices()
	_ = inv.SetPort()
	_ = inv.SetTimeout()
	_ = inv.CreateConnections()
	return
}
