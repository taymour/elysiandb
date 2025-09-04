# ElysianDB

**ElysianDB** is a lightweight HTTP key-value store written in Go. It offers a simple REST interface for storing arbitrary binary values and aims to be easy to run and benchmark.

This project is primarily an academic experiment to learn Go and explore building a minimalistic database system.

## Requirements

- [Go](https://go.dev/) 1.23+
- [k6](https://k6.io/) for optional load testing

## Configuration

Server settings are defined in `elysian.yaml`:

```yaml
folder: /tmp/elysiandb
host: localhost
port: 8089
```

- `folder` – path on disk where data files are stored
- `host` – interface the HTTP server binds to
- `port` – port used by the HTTP server

## Building and Running

```bash
# build executable
go build

# run directly
./elysiandb
# or
go run elysiandb.go
```

## HTTP API

All endpoints are rooted at the configured host and port.

| Method | Path         | Description                                                       |
|--------|--------------|-------------------------------------------------------------------|
| GET    | `/health`    | Liveness probe                                                    |
| PUT    | `/kv/{key}`  | Store value bytes for `key`, returns `204`                        |
| GET    | `/kv/{key}`  | Retrieve value bytes for `key`                                    |
| DELETE | `/kv/{key}`  | Remove value for `key`, returns `204`                             |
| POST   | `/save`      | Force persist current store to disk (already done automatically)  |
| POST   | `/reset`     | Clear all data from the store                                     |

### Examples

```bash
# store a value
curl -X PUT http://localhost:8089/kv/foo -d 'bar'

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

[MIT License](LICENSE)
