package rdap

import (
    "fmt"
    "net"
    "regexp"
    "strconv"
    "strings"
)

var (
    domainRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
    asnRegex = regexp.MustCompile(`^(?:AS)?(\d+)$`)
)

// ValidateDomain checks if a domain name is valid
func ValidateDomain(domain string) error {
    if domain == "" {
        return fmt.Errorf("domain cannot be empty")
    }
    if len(domain) > 255 {
        return fmt.Errorf("domain name too long")
    }
    if !domainRegex.MatchString(domain) {
        return fmt.Errorf("invalid domain format")
    }
    return nil
}

// ValidateIP checks if an IP address is valid
func ValidateIP(ip string) error {
    if ip == "" {
        return fmt.Errorf("IP cannot be empty")
    }
    parsedIP := net.ParseIP(ip)
    if parsedIP == nil {
        return fmt.Errorf("invalid IP address format")
    }
    return nil
}

// ValidateASN checks if an ASN is valid
func ValidateASN(asn string) error {
    if asn == "" {
        return fmt.Errorf("ASN cannot be empty")
    }

    // Remove "AS" prefix if present
    asn = strings.TrimPrefix(strings.ToUpper(asn), "AS")

    // Check if it matches the ASN format
    if !asnRegex.MatchString("AS" + asn) {
        return fmt.Errorf("invalid ASN format")
    }

    // Convert to number and check range
    num, err := strconv.ParseUint(asn, 10, 32)
    if err != nil {
        return fmt.Errorf("invalid ASN number: %w", err)
    }

    // Check if ASN is in valid range (0-4294967295)
    if num > 4294967295 {
        return fmt.Errorf("ASN number out of range (must be between 0 and 4294967295)")
    }

    return nil
}
