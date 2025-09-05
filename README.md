# ElysianDB

**ElysianDB** is a lightweight key–value store written in Go. It can serve **over TCP and/or HTTP**, with a minimal text protocol for TCP (à la Redis) and a tiny REST interface for HTTP. The project is primarily an academic experiment to learn Go and explore building a minimalistic database system.

---

## Requirements

* [Go](https://go.dev/) 1.23+
* Optional benchmarking:

  * **TCP**: built‑in `elysian_bench` load generator
  * **HTTP**: [k6](https://k6.io/)

---

## Configuration

Server settings are defined in `elysian.yaml`. You can enable **HTTP**, **TCP**, or **both**, and control the number of in‑memory **shards**.

```yaml
store:
  folder: /tmp/elysiandb        # on-disk data folder
  shards: 128                   # number of shards for the in-memory map (power of two recommended)

server:
  http:
    enabled: true
    host: 0.0.0.0
    port: 8089
  tcp:
    enabled: true
    host: 0.0.0.0
    port: 8088
```

**Keys**

* `store.folder` – path on disk where data files are stored
* `store.shards` – number of shards used for the in‑memory store (e.g., 128/256/512)
* `server.http.*` – HTTP listener configuration
* `server.tcp.*` – TCP listener configuration

> If only one protocol is desired, set the other one to `enabled: false`.

---

## Building and Running

```bash
# build executable
go build

# run directly with your local elysian.yaml
./elysiandb
# or
go run elysiandb.go
```

### Health

* HTTP: `GET /health` → `200 OK`
* TCP: send `PING` → `PONG` (when TCP is enabled)

---

## Docker

Prebuilt images are available on Docker Hub:

```bash
# pull the latest
docker pull taymour/elysiandb:latest

# or a specific version
docker pull taymour/elysiandb:v0.1.0
```

### Quick start (ephemeral)

```bash
docker run --rm \
  -p 8089:8089 \  # HTTP
  -p 8088:8088 \  # TCP
  taymour/elysiandb:latest
```

### With persistence

```bash
docker run -d --name elysiandb \
  -p 8089:8089 \
  -p 8088:8088 \
  -v elysian-data:/data \
  taymour/elysiandb:latest
```

The image includes a default config at `/app/elysian.yaml`:

```yaml
store:
  folder: /data
  shards: 128
server:
  http: { enabled: true,  host: 0.0.0.0, port: 8089 }
  tcp:  { enabled: true,  host: 0.0.0.0, port: 8088 }
```

To override the config, mount your own file:

```bash
docker run -d --name elysiandb \
  -p 8089:8089 -p 8088:8088 \
  -v elysian-data:/data \
  -v $(pwd)/elysian.yaml:/app/elysian.yaml:ro \
  taymour/elysiandb:latest
```

---

## Protocols

### TCP text protocol ("à la Redis")

A very small, line‑based text protocol. Each command is a line terminated by `\n`. Whitespace separates tokens.

**Supported commands (core):**

* `GET <key>` → returns raw value bytes; if missing, returns an empty payload or a not‑found marker
* `SET <key> <value>` → stores value; optional `TTL=<seconds>` support via `SET <key> <value> TTL=10`
* `DEL <key>` → deletes key
* `PING` → `PONG`

**Examples (netcat):**

```bash
# set a value
printf "SET foo bar\n" | nc 127.0.0.1 8088

# get it
printf "GET foo\n" | nc 127.0.0.1 8088

# set with ttl=10s
printf "SET session:123 abc TTL=10\n" | nc 127.0.0.1 8088

# delete
printf "DEL foo\n" | nc 127.0.0.1 8088
```

> The protocol is intentionally simple for benchmarking and learning purposes.

### HTTP API

| Method | Path                | Description                                                             |
| ------ | ------------------- | ----------------------------------------------------------------------- |
| GET    | `/health`           | Liveness probe                                                          |
| PUT    | `/kv/{key}?ttl=100` | Store value bytes for `key` with optional ttl in seconds, returns `204` |
| GET    | `/kv/{key}`         | Retrieve value bytes for `key`                                          |
| DELETE | `/kv/{key}`         | Remove value for `key`, returns `204`                                   |
| POST   | `/save`             | Force persist current store to disk (already done automatically)        |
| POST   | `/reset`            | Clear all data from the store                                           |

**Examples:**

```bash
# store a value
curl -X PUT http://localhost:8089/kv/foo -d 'bar'

# store a value for 10 seconds
curl -X PUT "http://localhost:8089/kv/foo?ttl=10" -d 'bar'

# fetch it
curl http://localhost:8089/kv/foo

# delete it
curl -X DELETE http://localhost:8089/kv/foo

# persist to disk
curl -X POST http://localhost:8089/save

# reset the store
curl -X POST http://localhost:8089/reset
```

---

## Benchmarks (local, indicative)

Two separate load generators are provided:

* **TCP**: native benchmark tool (`make tcp_benchmark`) (fully made by AI)
* **HTTP**: k6 script (`make http_benchmark`)

### TCP (paired SET→GET, payload 16B)

```
VUs: 500, duration: 20s, keys: 20,000, payload: 16B, paired mode (SET→GET 1:1)
Throughput: ~361,966 req/s  (≈180,983 pairs/s)
Latency:    p50 0.84 ms, p95 3.80 ms, p99 6.11 ms, max 28.53 ms
Errors:     0, Misses: 0
```

### HTTP (PUT/GET/DEL mix)

```
VUs: 200, duration: 30s
Requests:   ~72,814 req/s
Latency:    p50 ≈ 2.1 ms, p95 6.54 ms, max 45.83 ms
Errors:     0.00%
```

> The TCP benchmark runs a minimal text protocol with paired SET→GET and tiny payloads; the HTTP test mixes PUT/GET/DEL and pays the HTTP parsing overhead. As a result, **TCP typically delivers \~5× higher RPS** in this setup.

**Shortcuts:**

```bash
make tcp_benchmark   # runs the TCP benchmark tool with sensible defaults (fully made by AI)
make http_benchmark  # runs the k6 script (requires k6)
```

You can tweak VUs/duration/keys in scripts or via env vars as documented in the benchmark files.

---

## Testing

There are currently no unit tests.

---

## Debugging

If you experience boot issues, try removing `elysian*` data files in your configured `store.folder`. The project is evolving and may not be retro‑compatible.

---

[MIT License](LICENSE)
