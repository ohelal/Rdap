package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// IANA RDAP Bootstrap Service URLs
	ianaIPBootstrap    = "https://data.iana.org/rdap/ipv4.json"
	ianaASNBootstrap   = "https://data.iana.org/rdap/asn.json"
	ianaDNSBootstrap   = "https://data.iana.org/rdap/dns.json"
)

type RDAPClient struct {
	httpClient *http.Client
}

func NewRDAPClient() *RDAPClient {
	return &RDAPClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type bootstrapResponse struct {
	Services [][]interface{} `json:"services"`
}

func (c *RDAPClient) getBootstrapInfo(ctx context.Context, bootstrapURL string) (*bootstrapResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bootstrapURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating bootstrap request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching bootstrap info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from bootstrap server: %d", resp.StatusCode)
	}

	var bootstrap bootstrapResponse
	if err := json.NewDecoder(resp.Body).Decode(&bootstrap); err != nil {
		return nil, fmt.Errorf("decoding bootstrap response: %w", err)
	}

	return &bootstrap, nil
}

func (c *RDAPClient) findAppropriateServer(services [][]interface{}, query string) (string, error) {
	for _, service := range services {
		if len(service) < 2 {
			continue
		}

		ranges, ok := service[0].([]interface{})
		if !ok {
			continue
		}

		servers, ok := service[1].([]interface{})
		if !ok || len(servers) == 0 {
			continue
		}

		for _, r := range ranges {
			if strings.Contains(query, r.(string)) {
				serverURL := servers[0].(string)
				if !strings.HasPrefix(serverURL, "http") {
					serverURL = "https://" + serverURL
				}
				return serverURL, nil
			}
		}
	}
	return "", fmt.Errorf("no appropriate RDAP server found for query: %s", query)
}

func (c *RDAPClient) QueryIP(ctx context.Context, ip string) (map[string]interface{}, error) {
	bootstrap, err := c.getBootstrapInfo(ctx, ianaIPBootstrap)
	if err != nil {
		return nil, err
	}

	serverBase, err := c.findAppropriateServer(bootstrap.Services, ip)
	if err != nil {
		return nil, err
	}

	queryURL := fmt.Sprintf("%s/ip/%s", serverBase, url.PathEscape(ip))
	return c.makeRequest(ctx, queryURL)
}

func (c *RDAPClient) QueryASN(ctx context.Context, asn string) (map[string]interface{}, error) {
	bootstrap, err := c.getBootstrapInfo(ctx, ianaASNBootstrap)
	if err != nil {
		return nil, err
	}

	serverBase, err := c.findAppropriateServer(bootstrap.Services, asn)
	if err != nil {
		return nil, err
	}

	queryURL := fmt.Sprintf("%s/autnum/%s", serverBase, url.PathEscape(asn))
	return c.makeRequest(ctx, queryURL)
}

func (c *RDAPClient) QueryDomain(ctx context.Context, domain string) (map[string]interface{}, error) {
	bootstrap, err := c.getBootstrapInfo(ctx, ianaDNSBootstrap)
	if err != nil {
		return nil, err
	}

	serverBase, err := c.findAppropriateServer(bootstrap.Services, domain)
	if err != nil {
		return nil, err
	}

	queryURL := fmt.Sprintf("%s/domain/%s", serverBase, url.PathEscape(domain))
	return c.makeRequest(ctx, queryURL)
}

func (c *RDAPClient) makeRequest(ctx context.Context, url string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/rdap+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return result, nil
}
