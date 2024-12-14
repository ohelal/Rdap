// Package rdap provides a client for interacting with RDAP services
package rdap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Client represents an RDAP client
type Client struct {
	httpClient *http.Client
	baseURL    string
	metrics    *Metrics
	options    ClientOptions
}

// ClientOptions contains configuration options for the client
type ClientOptions struct {
	BaseURL       string
	Timeout       time.Duration
	Headers       map[string]string
	EnableMetrics bool
	RateLimit     int
	RetryOnLimit  bool
}

// NewClient creates a new RDAP client
func NewClient(options ...Option) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "http://localhost:8080",
		options: ClientOptions{},
	}

	// Apply options
	for _, opt := range options {
		opt(client)
	}

	return client
}

// Option defines a function type for configuring the client
type Option func(*Client)

// WithBaseURL sets the base URL for the RDAP service
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithMetrics enables metrics collection
func WithMetrics(enable bool) Option {
	return func(c *Client) {
		c.options.EnableMetrics = enable
		if enable {
			c.metrics = NewMetrics()
		}
	}
}

// WithHeaders sets custom HTTP headers
func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		c.options.Headers = headers
	}
}

// QueryDomain queries information about a domain
func (c *Client) QueryDomain(ctx context.Context, domain string) (map[string]interface{}, error) {
	if err := c.ValidateDomain(domain); err != nil {
		return nil, fmt.Errorf("invalid domain: %w", err)
	}

	url := fmt.Sprintf("%s/domain/%s", c.baseURL, domain)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return result, nil
}

// QueryIP queries information about an IP address
func (c *Client) QueryIP(ctx context.Context, ip string) (map[string]interface{}, error) {
	if err := c.ValidateIP(ip); err != nil {
		return nil, fmt.Errorf("invalid IP: %w", err)
	}

	url := fmt.Sprintf("%s/ip/%s", c.baseURL, ip)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return result, nil
}

// QueryASN queries information about an Autonomous System Number
func (c *Client) QueryASN(ctx context.Context, asn string) (map[string]interface{}, error) {
	if err := c.ValidateASN(asn); err != nil {
		return nil, fmt.Errorf("invalid ASN: %w", err)
	}

	// Remove "AS" prefix if present and convert to string
	asn = strings.TrimPrefix(strings.ToUpper(asn), "AS")

	url := fmt.Sprintf("%s/autnum/%s", c.baseURL, asn)
	return c.makeRequest(ctx, url)
}

// BatchQueryASNs queries multiple ASNs concurrently
func (c *Client) BatchQueryASNs(ctx context.Context, asns []string) map[string]interface{} {
	results := make(map[string]interface{})
	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, asn := range asns {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()

			resp, err := c.QueryASN(ctx, a)

			mutex.Lock()
			results[a] = map[string]interface{}{"result": resp, "error": err}
			mutex.Unlock()
		}(asn)
	}

	wg.Wait()
	return results
}

// GetMetrics returns the current metrics
func (c *Client) GetMetrics() *Metrics {
	return c.metrics
}

func (c *Client) makeRequest(ctx context.Context, url string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return result, nil
}

// ValidateDomain checks if a string is a valid domain name
func (c *Client) ValidateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name cannot be empty")
	}
	if strings.Contains(domain, " ") {
		return fmt.Errorf("domain name cannot contain spaces")
	}
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("domain name must contain at least one dot")
	}
	return nil
}

// ValidateIP checks if a string is a valid IP address
func (c *Client) ValidateIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP address cannot be empty")
	}
	if strings.Contains(ip, " ") {
		return fmt.Errorf("IP address cannot contain spaces")
	}
	// Basic format check - could be enhanced with net.ParseIP
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return fmt.Errorf("invalid IPv4 address format")
	}
	return nil
}

// ValidateASN checks if a string is a valid ASN
func (c *Client) ValidateASN(asn string) error {
	if asn == "" {
		return fmt.Errorf("ASN cannot be empty")
	}
	if !strings.HasPrefix(strings.ToUpper(asn), "AS") {
		return fmt.Errorf("ASN must start with 'AS'")
	}
	return nil
}
