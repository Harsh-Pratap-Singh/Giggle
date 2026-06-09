# Initialisation

How to set up, run, and verify the Giggle PageRank service.

## 1. Prerequisites

- **Go 1.21+** — verify with `go version`
- No other dependencies; the project uses only the Go standard library.

## 2. Project layout

```
Search_Engine/
├── pagerank/         # Library package (algorithm + tests)
├── cmd/
│   ├── cli/          # Standalone CLI demo
│   └── server/       # HTTP microservice
├── go.mod
├── CLAUDE.md
├── initialisation.md
└── README.md
```

## 3. First-time setup

From the `Search_Engine` directory:

```bash
go mod tidy
go build ./...
```

`go mod tidy` is a no-op the first time (no external imports), but it confirms the module resolves cleanly. `go build ./...` compiles every package and surfaces any import errors before you try to run things.

## 4. Run the tests

```bash
go test ./...
```

Expected: all tests in `pagerank/pagerank_test.go` pass — eight tests covering correctness, sort order, dangling nodes, parameter validation, and convergence.

For verbose output:

```bash
go test -v ./pagerank
```

## 5. Run the CLI 

```bash
go run ./cmd/cli
```

Expected output:

```
Total URLs processed: 4

Top 4 PageRank Results:
==================================================
page-c                                   | 0.3xxxxxxx
page-a                                   | 0.2xxxxxxx
...

PageRank convergence demonstration:
After  1 iterations - Top page: ...
After  5 iterations - Top page: ...
```

Override the defaults:

```bash
go run ./cmd/cli -damping=0.90 -iterations=100 -top=5
```

## 6. Run the HTTP microservice

```bash
go run ./cmd/server
```

The service listens on `:8080` by default. Override with `-addr=:9000`.

On startup, the server seeds itself with the same demo graph as the CLI, so the GET endpoints work immediately.

### Endpoints

| Method | Path                        | Purpose                                               |
| ------ | --------------------------- | ----------------------------------------------------- |
| GET    | `/healthz`                  | Liveness probe.                                       |
| GET    | `/pagerank`                 | All cached scores (sorted descending).                |
| GET    | `/pagerank?url=<page>`      | Score for a single page; 404 if unknown.              |
| POST   | `/compute`                  | Recompute scores from a JSON link graph; replaces cache. |

### Example calls

```bash
# health
curl http://localhost:8080/healthz

# all scores
curl http://localhost:8080/pagerank

# single page
curl 'http://localhost:8080/pagerank?url=page-c'

# recompute with a custom graph
curl -X POST http://localhost:8080/compute \
  -H 'Content-Type: application/json' \
  -d '{
    "backlinks": {
      "page-a": ["page-b", "page-c"],
      "page-b": ["page-c"],
      "page-c": ["page-a", "page-d"],
      "page-d": ["page-b", "page-c", "page-a"]
    },
    "outlinks_count": {
      "page-a": 2, "page-b": 1, "page-c": 2, "page-d": 1
    },
    "damping": 0.85,
    "iterations": 50
  }'
```

On Windows PowerShell, replace the `\` line continuations with a backtick (`` ` ``) and use double quotes around the JSON body.

## 7. Build standalone binaries

```bash
go build -o bin/pagerank-cli    ./cmd/cli
go build -o bin/pagerank-server ./cmd/server
```

Then run `./bin/pagerank-server -addr=:8080`.

## 8. Integrating into the wider search pipeline

The HTTP service is the integration surface for the rest of the system described in [README.md](README.md):

- **Backlinks Cluster / Indexer** → `POST /compute` to refresh scores after a crawl pass.
- **Query Engine** → `GET /pagerank?url=<doc>` to fetch the authority signal that gets blended with TF-IDF at query time.
- **Scaling Service** → `GET /healthz` for liveness.

Persistence (MongoDB) is not yet wired up — scores live in memory. Add a storage adapter in `cmd/server/main.go` when that backend is ready.

## 9. Troubleshooting

- **`package ... is not in std`** — run from the `Search_Engine` directory so Go resolves the local module.
- **Port already in use** — pass `-addr=:<other-port>` to the server.
- **404 on `/pagerank?url=...`** — the URL was not in the most recent compute pass. Either include it in the next `POST /compute` payload or query without the `url` param to see what is loaded.
