package rdap

import (
	"context"
	"sync"
)

// BatchQueryResult represents a result from a batch query
type BatchQueryResult struct {
	Result interface{}
	Error  error
}

// BatchQueryDomains queries multiple domains concurrently
func (c *Client) BatchQueryDomains(ctx context.Context, domains []string) map[string]interface{} {
	results := make(map[string]interface{})
	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, domain := range domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()

			if err := ValidateDomain(d); err != nil {
				mutex.Lock()
				results[d] = map[string]interface{}{"error": err}
				mutex.Unlock()
				return
			}

			resp, err := c.QueryDomain(ctx, d)

			mutex.Lock()
			results[d] = map[string]interface{}{"result": resp, "error": err}
			mutex.Unlock()
		}(domain)
	}

	wg.Wait()
	return results
}

// BatchQueryIPs queries multiple IP addresses concurrently
func (c *Client) BatchQueryIPs(ctx context.Context, ips []string) map[string]interface{} {
	results := make(map[string]interface{})
	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, ip := range ips {
		wg.Add(1)
		go func(i string) {
			defer wg.Done()

			if err := ValidateIP(i); err != nil {
				mutex.Lock()
				results[i] = map[string]interface{}{"error": err}
				mutex.Unlock()
				return
			}

			resp, err := c.QueryIP(ctx, i)

			mutex.Lock()
			results[i] = map[string]interface{}{"result": resp, "error": err}
			mutex.Unlock()
		}(ip)
	}

	wg.Wait()
	return results
}
