// Copyright (C) 2024 Helal <mohamed@helal.me>
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// SPDX-License-Identifier: AGPL-3.0-or-later

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
