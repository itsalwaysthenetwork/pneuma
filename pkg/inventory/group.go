package inventory

import "time"

// Metadata has information for filtering devices
type Metadata struct {
	Labels     map[string]string `yaml:"labels"`
	Platform   string            `yaml:"platform"` // Platform of the Device
	Groups     []string          `yaml:"groups"`
	Connection *struct {
		Port     int           `yaml:"port"`     // Port to connect to
		Username string        `yaml:"username"` // Username to use for authentication purposes
		Password string        `yaml:"password"` // Password to use for authentication purposes
		Timeout  time.Duration `yaml:"timeout"`  // Timeout to use for initial connection and all send commands
	} `yaml:"connection"`
}

// Group represents a network element/device
type Group struct {
	Name     string                 `yaml:"name"`     // Name is the name of the Device
	Spec     map[string]interface{} `yaml:"spec"`     // Spec defines a Device's attributes
	Metadata *Metadata              `yaml:"metadata"` // Metadata for connections and built-in filtering
}
