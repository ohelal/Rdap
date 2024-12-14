#!/bin/bash

# System tuning for 100k requests/second
echo "Tuning system for high-performance RDAP service..."

# File descriptor limits
ulimit -n 1000000
echo "* soft nofile 1000000" >> /etc/security/limits.conf
echo "* hard nofile 1000000" >> /etc/security/limits.conf

# Network tuning
cat >> /etc/sysctl.conf << EOF
# TCP tuning
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 262144
net.ipv4.tcp_max_syn_backlog = 262144
net.ipv4.tcp_fin_timeout = 5
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_max_tw_buckets = 2000000

# Memory tuning
net.core.rmem_max = 67108864
net.core.wmem_max = 67108864
net.ipv4.tcp_rmem = 4096 87380 67108864
net.ipv4.tcp_wmem = 4096 87380 67108864

# Connection tracking
net.netfilter.nf_conntrack_max = 1000000
net.ipv4.netfilter.ip_conntrack_max = 1000000

# Enable BBR
net.core.default_qdisc = fq
net.ipv4.tcp_congestion_control = bbr
EOF

sysctl -p
