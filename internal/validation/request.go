package validation

import (
	"net"
	"regexp"
	"strings"
)

type RequestValidator struct {
	ipv4Regex   *regexp.Regexp
	ipv6Regex   *regexp.Regexp
	asnRegex    *regexp.Regexp
	domainRegex *regexp.Regexp
}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{
		ipv4Regex:   regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`),
		ipv6Regex:   regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`),
		asnRegex:    regexp.MustCompile(`^AS\d+$`),
		domainRegex: regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`),
	}
}

func (v *RequestValidator) ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func (v *RequestValidator) ValidateASN(asn string) bool {
	return v.asnRegex.MatchString(strings.ToUpper(asn))
}

func (v *RequestValidator) ValidateDomain(domain string) bool {
	return v.domainRegex.MatchString(domain)
}
