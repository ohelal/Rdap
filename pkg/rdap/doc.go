// Copyright (C) 2024 Helal <mohamed@helal.me>
// SPDX-License-Identifier: AGPL-3.0-or-later

/*
Package rdap provides a Go client library for RDAP (Registration Data Access Protocol) lookups.
It supports IP address, ASN, and domain name queries with built-in caching and rate limiting.

# Basic Usage

Create a client and perform lookups:

	client := rdap.NewClient()

	// IP Lookup
	ipResult, err := client.LookupIP(context.Background(), "8.8.8.8")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("IP Owner: %s\n", ipResult.Name)

	// ASN Lookup
	asnResult, err := client.LookupASN(context.Background(), "15169")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ASN Name: %s\n", asnResult.Name)

	// Domain Lookup
	domainResult, err := client.LookupDomain(context.Background(), "example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Domain Status: %v\n", domainResult.Status)

# Custom Configuration

Create a client with custom configuration:

	client := rdap.NewClient(rdap.Config{
		BaseURL: "https://rdap.example.com",
		Timeout: 30 * time.Second,
		Cache: rdap.CacheConfig{
			Enabled: true,
			TTL:     time.Hour,
		},
		RateLimit: rdap.RateLimitConfig{
			RequestsPerSecond: 100,
			Burst:            10,
		},
	})

# Error Handling

The library returns specific error types that can be checked:

	result, err := client.LookupIP(ctx, "8.8.8.8")
	if err != nil {
		switch {
		case errors.Is(err, rdap.ErrNotFound):
			log.Printf("IP not found")
		case errors.Is(err, rdap.ErrRateLimit):
			log.Printf("Rate limit exceeded")
		case errors.Is(err, rdap.ErrTimeout):
			log.Printf("Request timed out")
		default:
			log.Printf("Unknown error: %v", err)
		}
	}

# Thread Safety

The client is safe for concurrent use by multiple goroutines.
*/
package rdap
