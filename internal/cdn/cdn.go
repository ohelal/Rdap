package cdn

// CDNConfig represents the configuration for CDN
type CDNConfig struct {
	Endpoint string
	Timeout  int
}

// NewCDNConfig creates a new CDN configuration
func NewCDNConfig(endpoint string, timeout int) *CDNConfig {
	return &CDNConfig{
		Endpoint: endpoint,
		Timeout:  timeout,
	}
}

// CDN represents a Content Delivery Network interface
type CDN interface {
	Get(key string) (interface{}, error)
	Put(key string, value interface{}) error
}

// NewCDN creates a new CDN instance
func NewCDN(config *CDNConfig) CDN {
	return &defaultCDN{
		config: config,
	}
}

type defaultCDN struct {
	config *CDNConfig
}

func (d *defaultCDN) Get(key string) (interface{}, error) {
	// Implement your CDN get logic here
	return nil, nil
}

func (d *defaultCDN) Put(key string, value interface{}) error {
	// Implement your CDN put logic here
	return nil
}
