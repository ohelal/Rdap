# API Documentation

This document describes the RDAP (Registration Data Access Protocol) REST API endpoints provided by this service. All responses follow the [RDAP JSON Response Format](https://tools.ietf.org/html/rfc7483).

## Base URL

```
http://localhost:8080
```

For production deployments, replace with your domain.

## Authentication

Currently, the API is unauthenticated. Rate limiting is applied based on client IP address.

## Common Headers

All responses include the following headers:

| Header | Description |
|--------|-------------|
| `Content-Type` | Always `application/rdap+json` |
| `X-Rate-Limit-Limit` | Maximum requests per hour |
| `X-Rate-Limit-Remaining` | Remaining requests in the current window |
| `X-Rate-Limit-Reset` | Time when the rate limit resets (Unix timestamp) |

## Endpoints

### IP Address Lookup

```http
GET /ip/{ip}
```

Lookup information about an IP address (IPv4 or IPv6).

**Parameters:**
- `ip` (path): IP address to lookup (e.g., "8.8.8.8" or "2001:db8::1")

**Example Request:**
```bash
curl -H "Accept: application/rdap+json" http://localhost:8080/ip/8.8.8.8
```

**Example Response:**
```json
{
  "objectClassName": "ip network",
  "handle": "NET-8-8-8-0-1",
  "startAddress": "8.8.8.0",
  "endAddress": "8.8.8.255",
  "ipVersion": "v4",
  "name": "GOOGLE-IPV4",
  "type": "ALLOCATION",
  "country": "US",
  "entities": [
    {
      "objectClassName": "entity",
      "handle": "GOGL",
      "roles": ["registrant"],
      "vcardArray": ["vcard", [
        ["version", {}, "text", "4.0"],
        ["fn", {}, "text", "Google LLC"]
      ]]
    }
  ]
}
```

### Domain Lookup

```http
GET /domain/{domain}
```

Lookup information about a domain name.

**Parameters:**
- `domain` (path): Domain name to lookup (e.g., "example.com")

**Example Request:**
```bash
curl -H "Accept: application/rdap+json" http://localhost:8080/domain/google.com
```

**Example Response:**
```json
{
  "objectClassName": "domain",
  "handle": "2138514_DOMAIN_COM-VRSN",
  "ldhName": "GOOGLE.COM",
  "nameservers": [
    {
      "objectClassName": "nameserver",
      "ldhName": "NS1.GOOGLE.COM"
    }
  ],
  "entities": [
    {
      "objectClassName": "entity",
      "handle": "MMR-2383",
      "roles": ["registrar"],
      "vcardArray": ["vcard", [
        ["version", {}, "text", "4.0"],
        ["fn", {}, "text", "MarkMonitor Inc."]
      ]]
    }
  ]
}
```

### ASN Lookup

```http
GET /autnum/{asn}
```

Lookup information about an Autonomous System Number.

**Parameters:**
- `asn` (path): ASN to lookup (e.g., "15169" or "AS15169")

**Example Request:**
```bash
curl -H "Accept: application/rdap+json" http://localhost:8080/autnum/AS15169
```

**Example Response:**
```json
{
  "objectClassName": "autnum",
  "handle": "AS15169",
  "startAutnum": 15169,
  "endAutnum": 15169,
  "name": "GOOGLE",
  "type": "DIRECT ALLOCATION",
  "entities": [
    {
      "objectClassName": "entity",
      "handle": "GOGL",
      "roles": ["registrant"],
      "vcardArray": ["vcard", [
        ["version", {}, "text", "4.0"],
        ["fn", {}, "text", "Google LLC"]
      ]]
    }
  ]
}
```

## Error Responses

The API uses standard HTTP status codes and returns error details in the response body.

### Example Error Response

```json
{
  "errorCode": 404,
  "title": "Not Found",
  "description": "The requested domain was not found",
  "notices": [
    {
      "title": "Terms of Use",
      "description": ["Service subject to Terms of Use."],
      "links": [
        {
          "value": "https://example.com/terms",
          "rel": "terms-of-service",
          "type": "text/html"
        }
      ]
    }
  ]
}
```

### Common Error Codes

| Status Code | Description |
|------------|-------------|
| 400 | Bad Request - Invalid input format |
| 404 | Not Found - Resource doesn't exist |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |
| 503 | Service Unavailable - Upstream service error |

## Rate Limiting

The API implements rate limiting to ensure fair usage:

- 1000 requests per hour per IP address
- Bursts of up to 100 requests are allowed
- Rate limits are applied globally across all endpoints

When a rate limit is exceeded, the API returns a 429 status code with a `Retry-After` header indicating when the client can resume making requests.

## Versioning

The API follows semantic versioning. The current version is included in the response headers:

```
X-API-Version: 1.0.0
```

## Additional Resources

- [RDAP Protocol Specification (RFC 7482)](https://tools.ietf.org/html/rfc7482)
- [RDAP Query Format (RFC 7482)](https://tools.ietf.org/html/rfc7482)
- [RDAP JSON Response Format (RFC 7483)](https://tools.ietf.org/html/rfc7483)
