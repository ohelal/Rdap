package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ohelal/rdap/pkg/rdap"
)

func main() {
	// Create client with metrics enabled
	client := rdap.NewClient(
		rdap.WithTimeout(10*time.Second),
		rdap.WithMetrics(true),
	)

	// Test single domain query
	fmt.Println("Testing single domain query...")
	domain, err := client.QueryDomain(context.Background(), "google.com")
	if err != nil {
		log.Printf("Error querying single domain: %v", err)
	} else {
		output, _ := json.MarshalIndent(domain, "", "    ")
		fmt.Printf("Domain Information:\n%s\n\n", output)
	}

	// Test batch domain queries
	fmt.Println("Testing batch domain queries...")
	domains := []string{"google.com", "github.com", "example.com"}
	results := client.BatchQueryDomains(context.Background(), domains)

	for domain, result := range results {
		resultMap := result.(map[string]interface{})
		if err, ok := resultMap["error"]; ok && err != nil {
			fmt.Printf("Error querying %s: %v\n", domain, err)
			continue
		}
		output, _ := json.MarshalIndent(resultMap["result"], "", "    ")
		fmt.Printf("Domain %s info:\n%s\n\n", domain, output)
	}

	// Test IP queries
	fmt.Println("Testing IP queries...")
	ips := []string{"8.8.8.8", "1.1.1.1"}
	ipResults := client.BatchQueryIPs(context.Background(), ips)

	for ip, result := range ipResults {
		resultMap := result.(map[string]interface{})
		if err, ok := resultMap["error"]; ok && err != nil {
			fmt.Printf("Error querying %s: %v\n", ip, err)
			continue
		}
		output, _ := json.MarshalIndent(resultMap["result"], "", "    ")
		fmt.Printf("IP %s info:\n%s\n\n", ip, output)
	}

	// Test ASN queries
	fmt.Println("\nTesting ASN queries...")
	asns := []string{"AS15169", "AS13335"} // Google and Cloudflare ASNs
	asnResults := client.BatchQueryASNs(context.Background(), asns)

	for asn, result := range asnResults {
		resultMap := result.(map[string]interface{})
		if err, ok := resultMap["error"]; ok && err != nil {
			fmt.Printf("Error querying %s: %v\n", asn, err)
			continue
		}
		output, _ := json.MarshalIndent(resultMap["result"], "", "    ")
		fmt.Printf("ASN %s info:\n%s\n\n", asn, output)
	}

	// Test ASN validation
	fmt.Println("Testing ASN validation...")
	if err := rdap.ValidateASN("invalid-asn"); err != nil {
		fmt.Printf("Validation caught invalid ASN: %v\n", err)
	}
	if err := rdap.ValidateASN("AS99999999999"); err != nil {
		fmt.Printf("Validation caught out-of-range ASN: %v\n", err)
	}

	// Test validation
	fmt.Println("Testing validation...")
	if err := rdap.ValidateDomain("invalid..domain"); err != nil {
		fmt.Printf("Validation caught invalid domain: %v\n", err)
	}
	if err := rdap.ValidateIP("invalid.ip"); err != nil {
		fmt.Printf("Validation caught invalid IP: %v\n", err)
	}

	// Print metrics
	if metrics := client.GetMetrics(); metrics != nil {
		fmt.Printf("\nMetrics:\n")
		fmt.Printf("Total Requests: %d\n", metrics.TotalRequests)
		fmt.Printf("Successful Requests: %d\n", metrics.SuccessfulRequests)
		fmt.Printf("Failed Requests: %d\n", metrics.FailedRequests)
		fmt.Printf("Average Latency: %v\n", metrics.GetAverageLatency())
	}
}
