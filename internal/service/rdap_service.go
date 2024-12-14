package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/ohelal/rdap/internal/config"
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
	"net"
	"sync"
)

// RDAPService represents the main service structure
type RDAPService struct {
	DNSConfig     *RDAPBootstrapConfig
	IPConfig      *RDAPBootstrapConfig
	ASNConfig     *RDAPBootstrapConfig
	ServiceConfig *config.Config
	client        *http.Client
	mu            sync.Mutex
}

// NewRDAPService creates a new RDAP service instance
func NewRDAPService(dnsConfig, ipConfig, asnConfig *RDAPBootstrapConfig, serviceConfig *config.Config) (*RDAPService, error) {
	if dnsConfig == nil || ipConfig == nil || asnConfig == nil {
		return nil, fmt.Errorf("all bootstrap configs must be non-nil")
	}
	if serviceConfig == nil {
		return nil, fmt.Errorf("service config must be non-nil")
	}

	return &RDAPService{
		DNSConfig:     dnsConfig,
		IPConfig:      ipConfig,
		ASNConfig:     asnConfig,
		ServiceConfig: serviceConfig,
		client: &http.Client{
			Timeout: serviceConfig.RDAP.Timeout,
		},
	}, nil
}

// findRDAPServerForASN finds the correct RDAP server for an ASN range
func (s *RDAPService) findRDAPServerForASN(asn int64) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for _, service := range s.ASNConfig.Services {
		if len(service) >= 2 {
			ranges, ok := service[0].([]interface{})
			if !ok {
				continue
			}
			servers, ok := service[1].([]interface{})
			if !ok || len(servers) == 0 {
				continue
			}
			
			for _, r := range ranges {
				asnRange, ok := r.(string)
				if !ok {
					continue
				}
				parts := strings.Split(asnRange, "-")
				if len(parts) == 2 {
					start, err1 := strconv.ParseInt(parts[0], 10, 64)
					end, err2 := strconv.ParseInt(parts[1], 10, 64)
					if err1 == nil && err2 == nil && asn >= start && asn <= end {
						if server, ok := servers[0].(string); ok {
							return server
						}
					}
				}
			}
		}
	}
	return ""
}

// findRDAPServerForIP finds the correct RDAP server for an IP range
func (s *RDAPService) findRDAPServerForIP(ipStr string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}

	for _, service := range s.IPConfig.Services {
		if len(service) >= 2 {
			cidrs, ok := service[0].([]interface{})
			if !ok {
				continue
			}
			servers, ok := service[1].([]interface{})
			if !ok || len(servers) == 0 {
				continue
			}
			
			for _, c := range cidrs {
				cidr, ok := c.(string)
				if !ok {
					continue
				}
				_, ipnet, err := net.ParseCIDR(cidr)
				if err == nil && ipnet.Contains(ip) {
					if server, ok := servers[0].(string); ok {
						return server
					}
				}
			}
		}
	}
	return ""
}

// findRDAPServerForTLD finds the correct RDAP server for a TLD
func (s *RDAPService) findRDAPServerForTLD(tld string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	tld = strings.ToLower(tld)
	for _, service := range s.DNSConfig.Services {
		if len(service) >= 2 {
			domains, ok := service[0].([]interface{})
			if !ok {
				continue
			}
			servers, ok := service[1].([]interface{})
			if !ok || len(servers) == 0 {
				continue
			}
			
			for _, d := range domains {
				domain, ok := d.(string)
				if !ok {
					continue
				}
				if strings.ToLower(domain) == tld {
					if server, ok := servers[0].(string); ok {
						return server
					}
				}
			}
		}
	}
	return ""
}

// forwardRequest handles the common logic for forwarding requests to RDAP servers
func (s *RDAPService) forwardRequest(c *fiber.Ctx, url string) error {
	resp, err := s.client.Get(url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"errorCode": 500,
			"title": "RDAP Server Error",
			"description": []string{err.Error()},
		})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"errorCode": 500,
			"title": "Response Error",
			"description": []string{err.Error()},
		})
	}

	c.Set("Content-Type", resp.Header.Get("Content-Type"))
	return c.Status(resp.StatusCode).Send(body)
}

// HandleIPLookup handles IP lookup requests
func (s *RDAPService) HandleIPLookup(c *fiber.Ctx) error {
	ip := c.Params("ip")
	rdapServer := s.findRDAPServerForIP(ip)
	if rdapServer == "" {
		return c.Status(404).JSON(fiber.Map{
			"errorCode": 404,
			"title": "IP Not Found",
			"description": []string{"No RDAP server found for IP: " + ip},
		})
	}

	return s.forwardRequest(c, rdapServer+"ip/"+ip)
}

// HandleDomainLookup handles domain lookup requests
func (s *RDAPService) HandleDomainLookup(c *fiber.Ctx) error {
	domain := strings.ToLower(c.Params("domain"))
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return c.Status(400).JSON(fiber.Map{
			"errorCode": 400,
			"title": "Invalid Domain",
			"description": []string{"Domain must include TLD"},
		})
	}
	
	tld := parts[len(parts)-1]
	rdapServer := s.findRDAPServerForTLD(tld)
	if rdapServer == "" {
		return c.Status(404).JSON(fiber.Map{
			"errorCode": 404,
			"title": "TLD Not Found",
			"description": []string{"No RDAP server found for TLD: " + tld},
		})
	}

	if !strings.HasSuffix(rdapServer, "/") {
		rdapServer += "/"
	}
	
	return s.forwardRequest(c, rdapServer+"domain/"+domain)
}

// HandleASNLookup handles ASN lookup requests
func (s *RDAPService) HandleASNLookup(c *fiber.Ctx) error {
	asnStr := c.Params("asn")
	asn, err := strconv.ParseInt(asnStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"errorCode": 400,
			"title": "Invalid ASN",
			"description": []string{"Invalid ASN format"},
		})
	}

	rdapServer := s.findRDAPServerForASN(asn)
	if rdapServer == "" {
		return c.Status(404).JSON(fiber.Map{
			"errorCode": 404,
			"title": "ASN Not Found",
			"description": []string{"No RDAP server found for ASN: " + asnStr},
		})
	}

	return s.forwardRequest(c, rdapServer+"autnum/"+asnStr)
}

// ReloadConfigs reloads all configurations
func (s *RDAPService) ReloadConfigs() error {
	// TODO: Implement config reloading logic
	return nil
} 