package cdn

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type CDNConfig struct {
	BaseURL    string
	Headers    http.Header
	MaxRetries int
	Timeout    int
}

type CDNManager struct {
	config CDNConfig
	client *http.Client
}

func NewCDNConfig(baseURL string) (*CDNConfig, error) {
	_, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	return &CDNConfig{
		BaseURL:    baseURL,
		Headers:    http.Header{},
		MaxRetries: 3,
		Timeout:    30,
	}, nil
}

func NewCDNManager(config CDNConfig) *CDNManager {
	return &CDNManager{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

func (c *CDNConfig) SetHeader(key, value string) {
	c.Headers.Set(key, value)
}

func (c *CDNConfig) GetHeader(key string) string {
	return c.Headers.Get(key)
}

func (cm *CDNManager) PurgeURL(url string) error {
	req, err := http.NewRequest("PURGE", url, nil)
	if err != nil {
		return err
	}

	req.Header = cm.config.Headers

	resp, err := cm.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to purge CDN cache: %d", resp.StatusCode)
	}

	return nil
}
