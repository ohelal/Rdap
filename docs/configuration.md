# Configuration Guide

This document describes how to configure the RDAP service using environment variables or configuration files.

## Environment Variables

### Server Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | HTTP server port | `8080` | No |
| `METRICS_PORT` | Prometheus metrics port | `9090` | No |
| `MAX_CONCURRENT_REQUESTS` | Maximum concurrent requests | `5000` | No |
| `TLS_CERT_FILE` | Path to TLS certificate | | No |
| `TLS_KEY_FILE` | Path to TLS private key | | No |

### Redis Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `REDIS_URL` | Redis connection URL | `redis:6379` | Yes |
| `REDIS_PASSWORD` | Redis password | | No |
| `REDIS_DB` | Redis database number | `0` | No |
| `CACHE_TTL` | Cache TTL in seconds | `3600` | No |

### Kafka Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `KAFKA_BROKERS` | Comma-separated Kafka brokers | | Yes |
| `KAFKA_TOPIC` | Kafka topic for events | `rdap-events` | No |
| `KAFKA_GROUP_ID` | Consumer group ID | `rdap-group` | No |
| `KAFKA_CLIENT_ID` | Client ID | `rdap-client` | No |

### Logging Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `LOG_LEVEL` | Log level (debug,info,warn,error) | `info` | No |
| `LOG_FORMAT` | Log format (json,text) | `json` | No |

### Rate Limiting
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `RATE_LIMIT_IP` | IP lookups per minute | `100` | No |
| `RATE_LIMIT_ASN` | ASN lookups per minute | `100` | No |
| `RATE_LIMIT_DOMAIN` | Domain lookups per minute | `100` | No |

## Configuration File

You can also use a YAML configuration file. Create `config.yaml`:

```yaml
server:
  port: 8080
  metrics_port: 9090
  max_concurrent_requests: 5000
  tls:
    cert_file: "/path/to/cert.pem"
    key_file: "/path/to/key.pem"

redis:
  url: "redis:6379"
  password: ""
  db: 0
  ttl: 3600

kafka:
  brokers:
    - "kafka-1:9092"
    - "kafka-2:9092"
  topic: "rdap-events"
  group_id: "rdap-group"
  client_id: "rdap-client"

logging:
  level: "info"
  format: "json"

rate_limits:
  ip: 100
  asn: 100
  domain: 100
```

## Using Configuration Files

1. Default locations checked:
   - `./config.yaml`
   - `./config/config.yaml`
   - `/etc/rdap/config.yaml`

2. Specify custom location:
   ```bash
   rdap --config /path/to/config.yaml
   ```

## Configuration Precedence

1. Command line flags (highest priority)
2. Environment variables
3. Configuration file
4. Default values (lowest priority)

## Kubernetes ConfigMap

For Kubernetes deployments, use a ConfigMap:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: rdap-config
  namespace: rdap
data:
  config.yaml: |
    server:
      port: 8080
    redis:
      url: "redis-master.rdap:6379"
    kafka:
      brokers:
        - "kafka-0.kafka-headless.rdap:9092"
    logging:
      level: "info"
```

Mount the ConfigMap in your deployment:

```yaml
volumes:
  - name: config
    configMap:
      name: rdap-config
volumeMounts:
  - name: config
    mountPath: /app/config
    readOnly: true
```

## Security Considerations

1. Never commit sensitive values to version control
2. Use Kubernetes Secrets for sensitive data
3. Rotate credentials regularly
4. Use TLS in production
5. Set appropriate file permissions

## Monitoring Configuration

The service exposes Prometheus metrics at `/metrics` including:
- Request latencies
- Cache hit rates
- Error rates
- Resource usage

Configure Prometheus to scrape these metrics:

```yaml
scrape_configs:
  - job_name: 'rdap'
    static_configs:
      - targets: ['rdap:9090']
```

## Troubleshooting

1. Check logs for configuration errors:
   ```bash
   LOG_LEVEL=debug rdap
   ```

2. Verify configuration:
   ```bash
   rdap config validate
   ```

3. Common issues:
   - Redis connection failures
   - Invalid Kafka configuration
   - Permission issues with TLS files
   - Incorrect file paths
