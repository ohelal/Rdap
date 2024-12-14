package kafka

import "time"

type Message struct {
	Type      string    `json:"type"`  // "ip", "domain", or "asn"
	Query     string    `json:"query"` // The actual query (IP address, domain name, or ASN)
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`    // Source of the request (e.g., "api", "cli")
	CacheHit  bool      `json:"cache_hit"` // Whether the result was from cache
}
