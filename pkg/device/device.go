package device

import (
	"fmt"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/network"
	"pneuma/pkg/device/l2proto"
	"time"
)

type Connection struct {
	*network.Driver
}

// GetLldpNeighbors gets LLDP neighbors based on Connection.Driver
func (d *Device) GetLldpNeighbors() (r *base.Reponse, err error) {
	r, err := d.SendCommand(l2proto.LldpShowCmd[d.Metadata.Platform])
	return
}

func (d *Device) SendCommand(s string, o ...base.SendOption) (r *base.Response, err error) {
	r, err = d.Connection.SendCommand(s, o...)
	fmt.Println(err)
	return
}

// Metadata has information for filtering devices
type Metadata struct {
	Labels     map[string]string `yaml:"labels"`
	Platform   string            `yaml:"platform"` // Platform of the Device
	Groups     []string          `yaml:"groups"`
	Connection *struct {
		Port     int           `yaml:"port"`     // Port to connect to
		Hostname string        `yaml:"hostname"` // Hostname/FQDN/IP to connect to
		Username string        `yaml:"username"` // Username to use for authentication purposes
		Password string        `yaml:"password"` // Password to use for authentication purposes
		Timeout  time.Duration `yaml:"timeout"`  // Timeout to use for initial connection and all send commands
	} `yaml:"connection"`
}

// Device represents a network element/device
type Device struct {
	Name       string                 `yaml:"name"`     // Name is the name of the Device
	Spec       map[string]interface{} `yaml:"spec"`     // Spec defines a Device's attributes
	Metadata   *Metadata              `yaml:"metadata"` // Metadata for connections and built-in filtering
	Connection *Connection            // Connection is the scrapligo SSH connection
}
