# API Reference

## Endpoints

### IP Address Lookup
```
GET /ip/{ip}
```

Lookup information about an IP address (IPv4 or IPv6).

**Parameters:**
- `ip` (path): IP address to lookup

**Example Response:**
```json
{
  "objectClassName": "ip network",
  "handle": "NET-8-8-8-0-1",
  "startAddress": "8.8.8.0",
  "endAddress": "8.8.8.255",
  "ipVersion": "v4",
  "name": "GOOGLE",
  "type": "ALLOCATION",
  "country": "US",
  "entities": [
    {
      "objectClassName": "entity",
      "handle": "GOGL",
      "roles": ["registrant"],
      "vcardArray": ["vcard", [["version", {}, "text", "4.0"]]]
    }
  ]
}
```

### ASN Lookup
```
GET /autnum/{asn}
```

Lookup information about an Autonomous System Number.

**Parameters:**
- `asn` (path): ASN to lookup (with or without "AS" prefix)

**Example Response:**
```json
{
  "objectClassName": "autnum",
  "handle": "AS15169",
  "startAutnum": 15169,
  "endAutnum": 15169,
  "name": "GOOGLE",
  "type": "DIRECT ALLOCATION",
  "country": "US",
  "entities": [
    {
      "objectClassName": "entity",
      "handle": "GOGL",
      "roles": ["registrant"],
      "vcardArray": ["vcard", [["version", {}, "text", "4.0"]]]
    }
  ]
}
```

### Domain Lookup
```
GET /domain/{domain}
```

Lookup information about a domain name.

**Parameters:**
- `domain` (path): Domain name to lookup

**Example Response:**
```json
{
  "objectClassName": "domain",
  "handle": "2138514_DOMAIN_COM-VRSN",
  "ldhName": "google.com",
  "status": ["client delete prohibited", "client transfer prohibited", "client update prohibited"],
  "entities": [
    {
      "objectClassName": "entity",
      "handle": "MMR-2383",
      "roles": ["registrar"],
      "vcardArray": ["vcard", [["version", {}, "text", "4.0"]]]
    }
  ],
  "events": [
    {
      "eventAction": "registration",
      "eventDate": "1997-09-15T04:00:00Z"
    }
  ],
  "nameservers": [
    {
      "objectClassName": "nameserver",
      "ldhName": "ns1.google.com"
    }
  ]
}
```

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "errorCode": 400,
  "title": "Bad Request",
  "description": "Invalid input parameter"
}
```

### 404 Not Found
```json
{
  "errorCode": 404,
  "title": "Not Found",
  "description": "Requested resource not found"
}
```

### 429 Too Many Requests
```json
{
  "errorCode": 429,
  "title": "Too Many Requests",
  "description": "Rate limit exceeded"
}
```

### 500 Internal Server Error
```json
{
  "errorCode": 500,
  "title": "Internal Server Error",
  "description": "An unexpected error occurred"
}
```

## Rate Limits

- IP lookups: 100 requests per minute
- ASN lookups: 100 requests per minute
- Domain lookups: 100 requests per minute

Rate limits are applied per client IP address. The following headers are included in responses:
- `X-RateLimit-Limit`: Maximum requests per minute
- `X-RateLimit-Remaining`: Remaining requests in the current window
- `X-RateLimit-Reset`: Time when the rate limit will reset (Unix timestamp)

## Authentication

Currently, the API does not require authentication. Rate limits are applied based on client IP address.
