package xray

import (
	"bytes"
	"x-ui/util/json_util"
)

type InboundConfig struct {
	Listen         json_util.RawMessage `json:"listen"` // listen 不能为空字符串
	Port           int                  `json:"port"`
	Protocol       string               `json:"protocol"`
	Settings       json_util.RawMessage `json:"settings"`
	StreamSettings json_util.RawMessage `json:"streamSettings"`
	Tag            string               `json:"tag"`
	Sniffing       json_util.RawMessage `json:"sniffing"`
	
	// Backend address configuration for proxy forwarding
	BackendAddress string `json:"backendAddress,omitempty"`
	BackendPort    int    `json:"backendPort,omitempty"`
	EnableBackend  bool   `json:"enableBackend,omitempty"`
}

func (c *InboundConfig) Equals(other *InboundConfig) bool {
	if !bytes.Equal(c.Listen, other.Listen) {
		return false
	}
	if c.Port != other.Port {
		return false
	}
	if c.Protocol != other.Protocol {
		return false
	}
	if !bytes.Equal(c.Settings, other.Settings) {
		return false
	}
	if !bytes.Equal(c.StreamSettings, other.StreamSettings) {
		return false
	}
	if c.Tag != other.Tag {
		return false
	}
	if !bytes.Equal(c.Sniffing, other.Sniffing) {
		return false
	}
	// Compare backend address fields
	if c.BackendAddress != other.BackendAddress {
		return false
	}
	if c.BackendPort != other.BackendPort {
		return false
	}
	if c.EnableBackend != other.EnableBackend {
		return false
	}
	return true
}
