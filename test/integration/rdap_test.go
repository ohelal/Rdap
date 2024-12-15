package integration

import (
	"context"
	"testing"
	"time"

	"github.com/ohelal/rdap/internal/connectivity"
	"github.com/ohelal/rdap/pkg/rdap"
	"github.com/stretchr/testify/assert"
)

func TestConnectivity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	err := connectivity.RunTests()
	assert.NoError(t, err, "Connectivity tests should pass")
}

func TestSingleDomainQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := rdap.NewClient(
		rdap.WithTimeout(10*time.Second),
		rdap.WithMetrics(true),
	)

	resp, err := client.QueryDomain(context.Background(), "google.com")
	assert.NoError(t, err, "Should successfully query single domain")
	assert.NotNil(t, resp, "Domain response should not be nil")

	// Basic validation of domain response
	assert.Contains(t, resp, "handle", "Domain response should contain a handle")
}

func TestBatchDomainQueries(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
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
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := rdap.NewClient(
		rdap.WithTimeout(10*time.Second),
		rdap.WithMetrics(true),
	)

	resp, err := client.QueryIP(context.Background(), "8.8.8.8")
	assert.NoError(t, err, "Should successfully query IP")
	assert.NotNil(t, resp, "IP response should not be nil")

	assert.Contains(t, resp, "handle", "IP response should contain a handle")
}

func TestASNQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	client := rdap.NewClient(
		rdap.WithTimeout(10*time.Second),
		rdap.WithMetrics(true),
	)

	resp, err := client.QueryASN(context.Background(), "AS15169")
	assert.NoError(t, err, "Should successfully query ASN")
	assert.NotNil(t, resp, "ASN response should not be nil")

	assert.Contains(t, resp, "handle", "ASN response should contain a handle")
}
