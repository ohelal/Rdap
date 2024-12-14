package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ohelal/rdap/internal/connectivity"
	"github.com/ohelal/rdap/pkg/rdap"
	"github.com/stretchr/testify/assert"
)

func TestConnectivity(t *testing.T) {
	err := connectivity.RunTests()
	assert.NoError(t, err, "Connectivity tests should pass")
}

func TestSingleDomainQuery(t *testing.T) {
	client := rdap.NewClient(
		rdap.WithTimeout(10*time.Second),
		rdap.WithMetrics(true),
	)

	domain, err := client.QueryDomain(context.Background(), "google.com")
	assert.NoError(t, err, "Should successfully query single domain")
	assert.NotNil(t, domain, "Domain response should not be nil")

	// Basic validation of domain response
	domainMap, ok := domain.(map[string]interface{})
	assert.True(t, ok, "Domain response should be a map")
	assert.Contains(t, domainMap, "handle", "Domain response should contain a handle")
}

func TestBatchDomainQueries(t *testing.T) {
	client := rdap.NewClient(
		rdap.WithTimeout(10*time.Second),
		rdap.WithMetrics(true),
	)

	domains := []string{"google.com", "github.com", "example.com"}
	results := client.BatchQueryDomains(context.Background(), domains)

	assert.Len(t, results, len(domains), "Should get results for all domains")

	for domain, result := range results {
		resultMap := result.(map[string]interface{})
		if err, ok := resultMap["error"]; ok && err != nil {
			t.Logf("Error querying %s: %v", domain, err)
			continue
		}
		
		assert.Contains(t, resultMap, "result", "Result for %s should contain 'result' field", domain)
		assert.NotNil(t, resultMap["result"], "Result for %s should not be nil", domain)
	}
}

func TestIPQuery(t *testing.T) {
	client := rdap.NewClient(
		rdap.WithTimeout(10*time.Second),
		rdap.WithMetrics(true),
	)

	ip, err := client.QueryIP(context.Background(), "8.8.8.8")
	assert.NoError(t, err, "Should successfully query IP")
	assert.NotNil(t, ip, "IP response should not be nil")

	ipMap, ok := ip.(map[string]interface{})
	assert.True(t, ok, "IP response should be a map")
	assert.Contains(t, ipMap, "handle", "IP response should contain a handle")
}

func TestASNQuery(t *testing.T) {
	client := rdap.NewClient(
		rdap.WithTimeout(10*time.Second),
		rdap.WithMetrics(true),
	)

	asn, err := client.QueryASN(context.Background(), "AS15169")
	assert.NoError(t, err, "Should successfully query ASN")
	assert.NotNil(t, asn, "ASN response should not be nil")

	asnMap, ok := asn.(map[string]interface{})
	assert.True(t, ok, "ASN response should be a map")
	assert.Contains(t, asnMap, "handle", "ASN response should contain a handle")
}
