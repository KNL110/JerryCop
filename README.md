# JerryCop

A high-performance, concurrent reverse proxy and layer-7 load balancer built in Go.

Designed for low-latency request routing and traffic management, JerryCop combines a CPU-sharded worker architecture with prefix-based path routing, active load balancing, and concurrency limiting to protect upstream backend services from overload.

---

## Key Features

- **High-Throughput Sharded Dispatcher**: Partitions incoming requests across dedicated CPU core shards to minimize channel lock contention under heavy concurrency.
- **Traffic Routing & Longest Prefix Matching (LPM)**: Intelligently routes paths (e.g. `/`, `/api`, `/api/v1`) to designated upstream server groups.
- **Pluggable Load Balancing**:
  - `ROUND_ROBIN`: Atomic state-tracked rotation across backend servers.
  - `RANDOM`: Low-overhead uniform pseudo-random distribution.
- **Built-in Rate Limiting & Concurrency Throttling**: Semaphore-guarded connection thresholds (`MaxConcReq`) to prevent server saturation, thread pool exhaustion, and backpressure collapse.
- **Optimized Network I/O**:
  - Upstream HTTP keep-alive connection pooling via reusable transports.
  - Recycled 32 KB copy buffers via `sync.Pool` to minimize GC pressure during large data transfers.
  - Fast single-pass target URL construction with zero intermediate URL re-parsing.
  - RFC 7230-compliant hop-by-hop header sanitization (`Connection`, `Keep-Alive`, `Upgrade`, etc.).

---

## Architecture Overview

```
                      ┌──────────────────────┐
                      │  Incoming HTTP Req   │
                      └──────────┬───────────┘
                                 │
                                 ▼
                     ┌────────────────────────┐
                     │    limitMiddleware     │ (Semaphore / MaxConcReq)
                     └───────────┬────────────┘
                                 │
                                 ▼
                     ┌────────────────────────┐
                     │      HTTPListener      │ (id % numShards)
                     └───────────┬────────────┘
                                 │
             ┌───────────────────┼───────────────────┐
             ▼                   ▼                   ▼
     ┌───────────────┐   ┌───────────────┐   ┌───────────────┐
     │ Shard Chan 0  │   │ Shard Chan 1  │   │ Shard Chan N  │
     └───────┬───────┘   └───────┬───────┘   └───────┬───────┘
             │                   │                   │
             ▼                   ▼                   ▼
     ┌───────────────┐   ┌───────────────┐   ┌───────────────┐
     │ Worker Pool 0 │   │ Worker Pool 1 │   │ Worker Pool N │
     └───────┬───────┘   └───────┬───────┘   └───────┬───────┘
             │                   │                   │
             └───────────────────┼───────────────────┘
                                 │
                                 ▼
                     ┌────────────────────────┐
                     │   Router & LoadBal     │ (LPM + RoundRobin / Random)
                     └───────────┬────────────┘
                                 │
                                 ▼
                     ┌────────────────────────┐
                     │   Upstream Backends    │
                     └────────────────────────┘
```

---

## Configuration

JerryCop is configured via `config/config.yaml`:

```yaml
Proxy:
  ProxyPort: ":8080"
  NumWorkers: 3000
  MaxConcReq: 20000
  DistAlgo: ROUND_ROBIN # Options: ROUND_ROBIN, RANDOM

Servers:
  Server1:
    ServerPorts:
      - "8001"
      - "8002"
    Names:
      - backend1
      - backend2
    Routing:
      Prefix: /
      ProxyPass:
        - http://127.0.0.1:8001
        - http://127.0.0.1:8002
      ProxyHeaders:
        Host: localhost
        XRealIP: 127.0.0.1
        XForwardedFor: 127.0.0.1

  Server2:
    ServerPorts:
      - "8003"
    Names:
      - api_backend
    Routing:
      Prefix: /api
      ProxyPass:
        - http://127.0.0.1:8003
      ProxyHeaders:
        Host: localhost
        XRealIP: 127.0.0.1
        XForwardedFor: 127.0.0.1
```

---

## Getting Started

### Prerequisites
- Go 1.22+

### Running the Proxy

```bash
# Run directly
go run ./cmd/main.go

# Or using Makefile
make run
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Or using Makefile
make test
```

### Building Binary

```bash
# Build standalone binary into bin/
make build
```
