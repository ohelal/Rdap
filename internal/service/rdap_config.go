package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// RDAPBootstrapConfig represents the RDAP bootstrap file format
type RDAPBootstrapConfig struct {
	Description string          `json:"description"`
	Publication string          `json:"publication"`
	Services    [][]interface{} `json:"services"`
	Version     string         `json:"version,omitempty"`
}

// LoadBootstrapConfig loads an RDAP bootstrap configuration from a file
func LoadBootstrapConfig(configPath string) (*RDAPBootstrapConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read bootstrap config file: %v", err)
	}

	var config RDAPBootstrapConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse bootstrap config: %v", err)
	}

	return &config, nil
}

// MergeIPConfigs merges IPv4 and IPv6 configs into a single config
func MergeIPConfigs(ipv4, ipv6 *RDAPBootstrapConfig) *RDAPBootstrapConfig {
	if ipv4 == nil || ipv6 == nil {
		if ipv4 != nil {
			return ipv4
		}
		return ipv6
	}

	mergedConfig := &RDAPBootstrapConfig{
		Description: "Merged IPv4 and IPv6 RDAP bootstrap file",
		Publication: time.Now().UTC().Format(time.RFC3339),
		Services:    make([][]interface{}, 0),
		Version:     "1.0",
	}

	mergedConfig.Services = append(mergedConfig.Services, ipv4.Services...)
	mergedConfig.Services = append(mergedConfig.Services, ipv6.Services...)

	return mergedConfig
}

// LoadAllBootstrapConfigs loads all bootstrap configurations from the config directory
func LoadAllBootstrapConfigs(configDir string) (*RDAPBootstrapConfig, *RDAPBootstrapConfig, *RDAPBootstrapConfig, error) {
	// Load DNS config
	dnsConfig, err := LoadBootstrapConfig(fmt.Sprintf("%s/dns.json", configDir))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load DNS config: %v", err)
	}

	// Load IPv4 config
	ipv4Config, err := LoadBootstrapConfig(fmt.Sprintf("%s/ipv4.json", configDir))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load IPv4 config: %v", err)
	}

	// Load IPv6 config
	ipv6Config, err := LoadBootstrapConfig(fmt.Sprintf("%s/ipv6.json", configDir))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load IPv6 config: %v", err)
	}

	// Load ASN config
	asnConfig, err := LoadBootstrapConfig(fmt.Sprintf("%s/asn.json", configDir))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load ASN config: %v", err)
	}

	// Merge IPv4 and IPv6 configs
	ipConfig := MergeIPConfigs(ipv4Config, ipv6Config)

	return dnsConfig, ipConfig, asnConfig, nil
}