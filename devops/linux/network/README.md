# Test communication between two networks

Two networks: netA and netB, with 2 services and a third acting as a gateway

## Architecture

| Network | Subnet | Services |
|---------|--------|----------|
| netA | 10.0.0.0/24 | serviceA (10.0.0.2), gateway (10.0.0.3) |
| netB | 10.1.0.0/24 | serviceB (10.1.0.2), gateway (10.1.0.3) |

## Quick Start

```bash
# Start services
docker-compose up 

# Test connectivity
docker exec gateway ping 10.0.0.2
docker exec gateway ping 10.1.0.2

# A can not talk with B
docker exec serviceA ping 10.1.0.2
docker exec serviceB ping 10.0.0.2

# create the IP route
docker exec serviceA ip route add 10.1.0.0/24 via 10.0.0.3
docker exec serviceB ip route add 10.0.0.0/24 via 10.1.0.3

# enable the gateway to foward package, on this image is enabled by default
# docker exec gateway echo 0 > /proc/sys/net/ipv4/ip_forward

# test the ping again, should work now
```