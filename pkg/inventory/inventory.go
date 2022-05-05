package inventory

import (
	"time"
	"pneuma/pkg/device"

	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
)

// Inventory represents a collection of Device
type Inventory struct {
	Devices map[string]*device.Device // Devices represents a collection of device.Device
	Groups  map[string]*Group         // Groups represent a generic collection of Group
}

// InventoryPlugin is an interface for Inventory plugins
// Currently, nothing implements this interface, but it exists as a soft suggestion on how to write other Inventory structs.
type InventoryPlugin interface {
	NewInventory() (*Inventory, error)
}

// MergeGroupsAndDevices updates all device.Device according to supported Group params
func (i *Inventory) MergeGroupsAndDevices() (err error) {
	for deviceID, device := range i.Devices {
		for _, g := range device.Metadata.Groups {
			group := i.Groups[g]
			if device.Metadata.Connection.Port == 0 && group.Metadata.Connection.Port != 0 {
				device.Metadata.Connection.Port = group.Metadata.Connection.Port
			}
			if device.Metadata.Connection.Timeout == 0 && group.Metadata.Connection.Timeout != 0 {
				device.Metadata.Connection.Timeout = group.Metadata.Connection.Timeout
			}
			if device.Metadata.Connection.Username == "" && group.Metadata.Connection.Username != "" {
				device.Metadata.Connection.Username = group.Metadata.Connection.Username
			}
			if device.Metadata.Connection.Password == "" && group.Metadata.Connection.Password != "" {
				device.Metadata.Connection.Password = group.Metadata.Connection.Password
			}
			device.Spec = mergeMaps(group.Spec, device.Spec)
			device.Metadata.Labels = mergeLabels(group.Metadata.Labels, device.Metadata.Labels)
			i.Devices[deviceID] = device
		}
	}
	return
}

// SetPort sets Inventory.Devices[].Metadata.Connection.Port to 22 if it isn't set
func (i *Inventory) SetPort() (err error) {
	for _, d := range i.Devices {
		if d.Metadata.Connection.Port == 0 {
			d.Metadata.Connection.Port = 22
		}
	}
	return nil
}

// SetTimeout sets a default timeout of 5 seconds in Inventory.Devices[].Metadata.Connection.Timeout if it isn't set
func (i *Inventory) SetTimeout() (err error) {
	for _, d := range i.Devices {
		if d.Metadata.Connection.Timeout == 0 {
			d.Metadata.Connection.Timeout = time.Second * 10
		}
	}
	return nil
}

// CreateConnections sets up the Scrapligo connections on Device.Connection based on Inventory.Devices[].Metadata.Connection
func (i *Inventory) CreateConnections() (err error) {
	for _, d := range i.Devices {
		d.Connection = &device.Connection{}
		d.Connection.Driver, _ = core.NewCoreDriver(
			d.Metadata.Connection.Hostname,
			d.Metadata.Platform,
			base.WithAuthUsername(d.Metadata.Connection.Username),
			base.WithAuthPassword(d.Metadata.Connection.Password),
			base.WithPort(d.Metadata.Connection.Port),
			base.WithTimeoutSocket(d.Metadata.Connection.Timeout),
		)
	}
	return nil
}

// mergeMaps recursively merges two map[string]interface{} objects
// Its primary use is to inherit all of Group.Spec into a device.Device.Spec based on device.Device.Groups membership
func mergeMaps(original, updated map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(original))
	for k, v := range original {
		merged[k] = v
	}
	for k, v := range updated {
		if v, ok := v.(map[string]interface{}); ok {
			if updatedValue, ok := merged[k]; ok {
				if updatedValue, ok := updatedValue.(map[string]interface{}); ok {
					merged[k] = mergeMaps(updatedValue, v)
					continue
				}
			}
		}
		merged[k] = v
	}
	return merged
}

func mergeLabels(groupLabels, deviceLabels map[string]string) map[string]string {
	mergedLabels := make(map[string]string)
	for k, v := range groupLabels {
		mergedLabels[k] = v
	}
	for k, v := range deviceLabels {
		mergedLabels[k] = v
	}
	return mergedLabels
}
