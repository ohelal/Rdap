# Kubernetes Deployment Guide

This guide provides detailed instructions for deploying the RDAP service on Kubernetes.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Detailed Setup](#detailed-setup)
- [Configuration](#configuration)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Tools
- Kubernetes cluster (v1.22+)
- kubectl (latest version)
- Helm v3 (optional, for chart installation)
- Docker (for building custom images)

### Resource Requirements
- Minimum 2 CPU cores
- 4GB RAM
- 20GB storage

## Architecture

The RDAP service on Kubernetes consists of the following components:

```
                                    [Ingress Controller]
                                           │
                                           ▼
                                    [RDAP Service]
                                     │         │
                              ┌─────┘         └─────┐
                              ▼                     ▼
                        [Redis Cache]         [Kafka Cluster]
```

## Quick Start

1. **Clone the Repository**
   ```bash
   git clone https://github.com/ohelal/rdap.git
   cd rdap
   ```

2. **Apply Base Configuration**
   ```bash
   kubectl create namespace rdap
   kubectl apply -f k8s/namespace.yaml
   kubectl apply -f k8s/configmap.yaml
   kubectl apply -f k8s/secret.yaml
   ```

3. **Deploy Dependencies**
   ```bash
   # Deploy Redis
   kubectl apply -f k8s/redis/
   
   # Deploy Kafka
   kubectl apply -f k8s/kafka/
   ```

4. **Deploy RDAP Service**
   ```bash
   kubectl apply -f k8s/rdap/
   ```

## Detailed Setup

### 1. Namespace and RBAC

Create dedicated namespace and required permissions:

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: rdap

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: rdap-service
  namespace: rdap

---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: rdap-service-role
  namespace: rdap
rules:
  - apiGroups: [""]
    resources: ["configmaps", "secrets"]
    verbs: ["get", "list", "watch"]
```

### 2. ConfigMaps and Secrets

Store configuration:

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: rdap-config
  namespace: rdap
data:
  CACHE_TTL: "3600"
  MAX_CONCURRENT_REQUESTS: "5000"
  LOG_LEVEL: "info"
  METRICS_ENABLED: "true"

---
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: rdap-secret
  namespace: rdap
type: Opaque
data:
  REDIS_PASSWORD: <base64-encoded-password>
  KAFKA_USERNAME: <base64-encoded-username>
  KAFKA_PASSWORD: <base64-encoded-password>
```

### 3. Redis Deployment

```yaml
# k8s/redis/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  namespace: rdap
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:7.0-alpine
        ports:
        - containerPort: 6379
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "200m"
```

### 4. RDAP Service Deployment

```yaml
# k8s/rdap/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rdap-service
  namespace: rdap
spec:
  replicas: 3
  selector:
    matchLabels:
      app: rdap-service
  template:
    metadata:
      labels:
        app: rdap-service
    spec:
      containers:
      - name: rdap-service
        image: ohelal/rdap:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: rdap-config
        - secretRef:
            name: rdap-secret
        resources:
          requests:
            memory: "512Mi"
            cpu: "200m"
          limits:
            memory: "1Gi"
            cpu: "500m"
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 15
          periodSeconds: 20
```

### 5. Service and Ingress

```yaml
# k8s/rdap/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: rdap-service
  namespace: rdap
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: rdap-service

---
# k8s/rdap/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: rdap-ingress
  namespace: rdap
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: rdap.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: rdap-service
            port:
              number: 80
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| REDIS_URL | Redis connection string | redis:6379 |
| CACHE_TTL | Cache duration in seconds | 3600 |
| MAX_CONCURRENT_REQUESTS | Maximum concurrent requests | 5000 |
| LOG_LEVEL | Logging level | info |
| METRICS_ENABLED | Enable Prometheus metrics | true |

### Resource Scaling

Adjust replicas based on load:

```yaml
# k8s/rdap/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: rdap-hpa
  namespace: rdap
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: rdap-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

## Monitoring

### Prometheus Metrics

The service exposes metrics at `/metrics` endpoint. Configure Prometheus to scrape these metrics:

```yaml
# k8s/monitoring/servicemonitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: rdap-monitor
  namespace: rdap
spec:
  selector:
    matchLabels:
      app: rdap-service
  endpoints:
  - port: metrics
```

### Grafana Dashboard

Import the provided Grafana dashboard for monitoring:
- Request rates
- Response times
- Cache hit/miss ratios
- Error rates
- Resource usage

## Troubleshooting

### Common Issues

1. **Pod Startup Issues**
   ```bash
   kubectl get pods -n rdap
   kubectl describe pod <pod-name> -n rdap
   kubectl logs <pod-name> -n rdap
   ```

2. **Service Connectivity**
   ```bash
   kubectl get svc -n rdap
   kubectl get endpoints -n rdap
   ```

3. **Cache Issues**
   ```bash
   kubectl exec -it <redis-pod> -n rdap -- redis-cli
   > INFO
   > MONITOR
   ```

### Health Checks

Monitor service health:
```bash
kubectl get pods -n rdap -o wide
kubectl top pods -n rdap
```

## Security

1. **Network Policies**
   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: rdap-network-policy
     namespace: rdap
   spec:
     podSelector:
       matchLabels:
         app: rdap-service
     ingress:
     - from:
       - namespaceSelector:
           matchLabels:
             name: ingress-nginx
     egress:
     - to:
       - podSelector:
           matchLabels:
             app: redis
       - podSelector:
           matchLabels:
             app: kafka
   ```

2. **Pod Security**
   ```yaml
   securityContext:
     runAsNonRoot: true
     runAsUser: 1000
     allowPrivilegeEscalation: false
   ```

## Updates and Maintenance

### Rolling Updates
```bash
kubectl set image deployment/rdap-service rdap-service=ohelal/rdap:new-version -n rdap
```

### Backup and Restore
1. Redis data backup
2. Configuration backup
3. Restore procedures

For more information or support, please visit our [GitHub repository](https://github.com/ohelal/rdap) or contact mohamed@helal.me.
