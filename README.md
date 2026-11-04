# JerryCop

A high-performance, concurrent reverse proxy and layer-7 load balancer built in Go.

Designed for low-latency request routing and traffic management, JerryCop combines a CPU-sharded worker architecture with prefix-based path routing, active load balancing, and concurrency limiting to protect upstream backend services from overload.

---

## Key Features

- **High-Throughput Sharded Dispatcher**: Partitions incoming requests across dedicated CPU core shards (`runtime.NumCPU()`) to minimize channel lock contention under heavy concurrency.
- **Concurrent Worker Pools & Graceful Draining**: Workers pull jobs asynchronously from shard channels, executing requests in per-job goroutines tracked by `sync.WaitGroup` to enable clean in-flight request draining on shutdown.
- **Traffic Routing & Longest Prefix Matching (LPM)**: Intelligently routes paths (e.g. `/`, `/api`, `/api/v1`) to designated upstream server groups, supporting multiple prefix routes per server.
- **Pluggable Load Balancing**:
  - `ROUND_ROBIN`: Atomic state-tracked rotation across backend servers.
  - `RANDOM`: Low-overhead uniform pseudo-random distribution.
  - `LEAST_CONNECTION`: Upstream selection based on active connection tracking.
- **Rate Limiting & Lifecycle Synchronization**: Semaphore-guarded connection thresholds (`MaxConcReq`) coordinated with per-request `Done` channel signaling to prevent server saturation and cleanly release slots.
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
                     │      HTTPListener      │◄────────────────┐
                     └───────────┬────────────┘                 │
                                 │ (id % numShards)             │
             ┌───────────────────┼───────────────────┐          │
             ▼                   ▼                   ▼          │
     ┌───────────────┐   ┌───────────────┐   ┌───────────────┐  │ (close job.Done /
     │ Shard Chan 0  │   │ Shard Chan 1  │   │ Shard Chan N  │  │  release slot)
     └───────┬───────┘   └───────┬───────┘   └───────┬───────┘  │
             │                   │                   │          │
             ▼                   ▼                   ▼          │
     ┌───────────────┐   ┌───────────────┐   ┌───────────────┐  │
     │ Worker Pool 0 │   │ Worker Pool 1 │   │ Worker Pool N │  │ (Async goroutines
     └───────┬───────┘   └───────┬───────┘   └───────┬───────┘  │  + sync.WaitGroup)
             │                   │                   │          │
             └───────────────────┼───────────────────┘          │
                                 │                              │
                                 ▼                              │
                     ┌────────────────────────┐                 │
                     │   Router & LoadBal     │                 │
                     └───────────┬────────────┘                 │
                                 │ (LPM + RR/Random/LeastConn)  │
                                 ▼                              │
                     ┌────────────────────────┐                 │
                     │   Upstream Backends    │─────────────────┘
                     └────────────────────────┘ (stream response)
```

---

## Configuration

JerryCop is configured via `config/config.yaml`:

```yaml
Proxy:
  ProxyPort: ":8080"
  NumWorkers: 3000
  MaxConcReq: 20000
  DistAlgo: ROUND_ROBIN # Options: ROUND_ROBIN, RANDOM, LEAST_CONNECTION

Servers:
  Server1:
    ServerPorts:
      - "8001"
      - "8002"
      - "8003"
    Names:
      - backend1
      - backend2
    Routing:
      - Prefix: /api
        ProxyPass:
          - http://127.0.0.1:8003
        ProxyHeaders:
          Host: localhost
          XRealIP: 127.0.0.1
          XForwardedFor: 127.0.0.1
      - Prefix: /
        ProxyPass:
          - http://127.0.0.1:8001
          - http://127.0.0.1:8002
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
