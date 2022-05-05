// This is only an example and isn't an actual tool that you should use.
package main

import (
	"fmt"
	"pneuma/pkg/device"
	"pneuma/pkg/inventory"

	"gopkg.in/yaml.v3"
)

func main() {
	i := inventory.YamlInventory{Filename: "example.yml"}
	inv, err := i.NewInventory()

	// Send some commands
	inv.Devices["edgeRouter1"].Connection.Open()
	defer inv.Devices["edgeRouter1"].Connection.Close()
	r, _ := inv.Devices["edgeRouter1"].GetLldpNeighbors()
	fmt.Println(r)
	fmt.Println("---")
	r, _ = inv.Devices["edgeRouter1"].Connection.SendCommand("show version")
	fmt.Println(r)
	fmt.Println("---")

	// Filter the inventory
	newInv, _ := inv.FilterMetadata(
		&device.Metadata{
			Labels: map[string]string{
				"dc": "slc1",
			},
			Platform: "arista_eos",
		})
	fmt.Println(newInv.Devices)
	fmt.Println("---")
	y, _ := yaml.Marshal(newInv.Devices["edgeRouter1"].Spec)
	fmt.Println(string(y))
}
