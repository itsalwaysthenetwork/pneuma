package inventory

import (
	"errors"
	"pneuma/pkg/device"
	"regexp"
)

// CustomFilterFunc is a function that can be used to filter the inventory
type CustomFilterFunc func(*device.Device) (bool, error)

// Filter filters the hosts in the Inventory, returning a copy of the current
// Inventory instance but with only the hosts that passed the filter
func (i *Inventory) Filter(f CustomFilterFunc) (*Inventory, error) {
	filtered := &Inventory{
		Devices: make(map[string]*device.Device),
	}
	for deviceId, device := range i.Devices {
		match, err := f(device)
		if err != nil {
			return filtered, errors.New("custom filter function returned error")
		}
		if match {
			filtered.Devices[deviceId] = device
		}
	}
	return filtered, nil
}

// FilterNameRegex filters an Inventory's device.Device.Name by the regex s
func (i *Inventory) FilterNameRegex(s string) (*Inventory, error) {
	regex := regexp.MustCompile(s)
	return i.Filter(func(d *device.Device) (bool, error) {
		if regex.MatchString(d.Name) {
			return true, nil
		}
		return false, nil
	})
}

// FilterMetadata filters an inventory based on Inventory.Devices[].Metadata
// Currently, only supports device.Metadata.Labels and device.Metadata.Platform
func (i *Inventory) FilterMetadata(m *device.Metadata) (*Inventory, error) {
	if m.Connection != nil {
		return nil, errors.New("`m.Connection`: unsupported")
	}

	return i.Filter(func(d *device.Device) (bool, error) {
		if (m.Platform != "" && m.Platform != d.Metadata.Platform) || !matchesLabels(d, m.Labels) {
			return false, nil
		}
		return true, nil
	})
}

// matchesLabels checks for a subset of labels in Device.Metadata.Labels
func matchesLabels(d *device.Device, labels map[string]string) bool {
	for k, v := range labels {
		labelValue, hasLabel := d.Metadata.Labels[k]
		if !hasLabel || labelValue != v {
			return false
		}
	}
	return true
}
