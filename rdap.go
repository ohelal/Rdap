// Package rdap provides a high-performance Registration Data Access Protocol (RDAP)
// service and command-line tool for IP, ASN, and domain lookups.
//
// Features:
//   - High-performance RDAP lookups for IP addresses, ASNs, and domains
//   - Built-in distributed caching with Redis
//   - Message queue integration with Kafka
//   - CLI and Service modes
//   - Standards-compliant RDAP responses
//   - Kubernetes-ready deployment
//
// Example usage:
//
//	client := rdap.NewClient()
//	
//	// IP Lookup
//	ipResult, err := client.LookupIP(context.Background(), "8.8.8.8")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("IP Owner: %s\n", ipResult.Name)
//	
//	// ASN Lookup
//	asnResult, err := client.LookupASN(context.Background(), "15169")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("ASN Name: %s\n", asnResult.Name)
//	
//	// Domain Lookup
//	domainResult, err := client.LookupDomain(context.Background(), "example.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Domain Status: %v\n", domainResult.Status)
package rdap

// Version is the current version of the RDAP package
const Version = "0.1.1"
