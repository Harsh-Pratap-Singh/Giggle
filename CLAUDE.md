# CLAUDE.md

Guidance for Claude Code when working in this repository.

## Project

**Giggle PageRank Service** — Go implementation of the PageRank algorithm, packaged as a microservice for a larger distributed search engine pipeline (crawler → indexer → TF-IDF → **PageRank** → query engine).

The owning user is `Harsh-Pratap-Singh`. The wider system context is documented in [README.md](README.md); this repo only owns the PageRank component.

## Layout

```
Search_Engine/
├── pagerank/              # Core algorithm package (importable, no main)
│   ├── pagerank.go        # PageRankCalculator + Calculate()
│   └── pagerank_test.go   # Unit tests + convergence checks
├── cmd/
│   ├── cli/main.go        # Standalone CLI demo (was the original PageRank.go main)
│   └── server/main.go     # HTTP microservice exposing /pagerank, /healthz
├── go.mod
├── CLAUDE.md              # This file
├── initialisation.md      # Setup + run instructions
└── README.md              # Algorithm + system architecture docs
```

The original flat `PageRank.go` was split into `pagerank/pagerank.go` (library) and `cmd/cli/main.go` (demo). The library has no `main` so it can be imported by the HTTP server and other consumers.

## Conventions

- **Module path**: `github.com/Harsh-Pratap-Singh/Search_Engine`. Imports inside the repo use this prefix.
- **Standard library only**: no third-party deps. Keep it that way unless the user explicitly approves a dependency — the README sells "no external dependencies" as a feature.
- **Damping default 0.85, iterations 50**: matches Google's published values. Don't change defaults without a reason.
- **Builder pattern**: `NewPageRankCalculator().SetDamping(d).SetIterations(n)` — keep new options as chainable setters that validate and return the receiver.
- **Map-based graph input**: `backlinks map[string][]string` and `outlinksCount map[string]int`. Don't introduce a new graph type without discussing — downstream services (Backlinks Cluster, Indexer) feed these shapes directly.

## Running

See [initialisation.md](initialisation.md) for full instructions. Quick reference:

```bash
go test ./...                              # tests
go run ./cmd/cli                           # demo run
go run ./cmd/server                        # HTTP service on :8080
curl 'http://localhost:8080/pagerank?url=page-a'
```

## When making changes

- Algorithm changes (damping math, dangling-node handling, convergence) must keep the existing tests green and add new ones — PageRank's correctness is load-bearing for the whole search pipeline.
- HTTP-shape changes to `/pagerank` are a breaking change for the Query Engine consumer. Flag these explicitly.
- Don't add CLAUDE.md-style planning docs, design files, or scratch markdown unless asked. README.md and initialisation.md are the only docs the user has approved.
