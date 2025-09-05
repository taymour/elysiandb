# ElysianDB

**ElysianDB** is a lightweight HTTP key-value store written in Go. It offers a simple REST interface for storing arbitrary binary values and aims to be easy to run and benchmark.

This project is primarily an academic experiment to learn Go and explore building a minimalistic database system.

## Requirements

* [Go](https://go.dev/) 1.23+
* [k6](https://k6.io/) for optional load testing

## Configuration

Server settings are defined in `elysian.yaml`:

```yaml
folder: /tmp/elysiandb
host: localhost
port: 8089
```

* `folder` – path on disk where data files are stored
* `host` – interface the HTTP server binds to
* `port` – port used by the HTTP server

## Building and Running

```bash
# build executable
go build

# run directly
./elysiandb
# or
go run elysiandb.go
```

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
docker run --rm -p 8089:8089 taymour/elysiandb:latest
curl -I -X GET http://localhost:8089/health
```

### With persistence

```bash
docker run -d --name elysiandb \
  -p 8089:8089 \
  -v elysian-data:/data \
  taymour/elysiandb:latest
```

The image includes a default config at `/app/elysian.yaml`:

```yaml
folder: /data
host: 0.0.0.0
port: 8089
```

To override the config, mount your own file:

```bash
docker run -d --name elysiandb \
  -p 8089:8089 \
  -v elysian-data:/data \
  -v $(pwd)/elysian.yaml:/app/elysian.yaml:ro \
  taymour/elysiandb:latest
```

**Healthcheck**: `GET /health`

## HTTP API

All endpoints are rooted at the configured host and port.

| Method | Path                | Description                                                             |
| ------ | ------------------- | ----------------------------------------------------------------------- |
| GET    | `/health`           | Liveness probe                                                          |
| PUT    | `/kv/{key}?ttl=100` | Store value bytes for `key` with optional ttl in seconds, returns `204` |
| GET    | `/kv/{key}`         | Retrieve value bytes for `key`                                          |
| DELETE | `/kv/{key}`         | Remove value for `key`, returns `204`                                   |
| POST   | `/save`             | Force persist current store to disk (already done automatically)        |
| POST   | `/reset`            | Clear all data from the store                                           |

### Examples

```bash
# store a value
curl -X PUT http://localhost:8089/kv/foo -d 'bar'

# store a value for 10 seconds
curl -X PUT http://localhost:8089/kv/foo?ttl=10 -d 'bar'

# fetch it
curl http://localhost:8089/kv/foo

# delete it
curl -X DELETE http://localhost:8089/kv/foo

# persist to disk
curl -X POST http://localhost:8089/save

# reset the store
curl -X POST http://localhost:8089/reset
```

## Testing

There are currently no unit tests.

## Benchmarking with k6

A k6 script is provided in `elysian_k6.js`. Run it with:

```bash
BASE_URL=http://localhost:8089 KEYS=5000 VUS=200 DURATION=30s k6 run elysian_k6.js
```

A sample local run yielded for 5000 keys in the store :

```
200 virtual users ran for 30, generating 1,953,090 requests (~65,100 req/s).

Average request latency was 2.75 ms with p95 at 6.78 ms.

Zero HTTP errors occurred; all checks succeeded.

Throughput was about 5.4 MB/s received and 11 MB/s sent.

Overall, the system handled the load with stable performance and no failures.
```

The same command is available in `benchmark.sh`.

## Debugging

If ever you experience problems with booting the DB, try removing the elysian\*.db files. The project is evolving and it may not be retro-compatible for the moment.

[MIT License](LICENSE)
