# Compete Finder

A Go service that finds competitor startups for a given product description by scoring real Y Combinator companies with TF-IDF, served over Connect RPC.

Given a name, description, and sector, it returns the top-N most similar YC startups (with match scores) and a market heatmap showing how that sector has grown across YC batches.

## Architecture

```
internal/provider/         Startup struct, Provider interface
internal/provider/yc/      YC API adapter (yc-oss.github.io)
internal/matcher/          Tokenizer + TF-IDF scoring
internal/heatmap/          Sector filtering + batch trend builder
proto/compete/v1/          Protobuf definitions
gen/compete/v1/            Generated Connect RPC code
internal/handler/          Connect service handlers
cmd/server/                Server entry point
```

## Requirements

- Go 1.22+

## Running the server

```bash
go run ./cmd/server
# listening on :8080
```

The server fetches the YC company list on demand from `https://yc-oss.github.io/api/companies/all.json`, so it needs network access on the first request.

## Calling the RPCs

Connect supports HTTP/JSON, so you can call the service with plain `curl`.

### FindCompetitors

Returns the top N YC startups in a given sector ranked by TF-IDF similarity to your description.

```bash
curl -X POST http://localhost:8080/compete.v1.CompeteService/FindCompetitors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MyFintech",
    "description": "fintech payments platform",
    "sector": "fintech",
    "limit": 3
  }'
```

### GetMarketHeatmap

Returns YC batch counts and per-batch trend (`growing` / `stable` / `shrinking`) for a sector, plus an overall market status and growth factor.

```bash
curl -X POST http://localhost:8080/compete.v1.CompeteService/GetMarketHeatmap \
  -H "Content-Type: application/json" \
  -d '{"sector": "fintech"}'
```

## Tests

```bash
go test ./...                        # all packages
go test ./internal/handler/...       # handler integration tests
```

The handler tests stand up an in-process `httptest.Server` with a fake provider, so they run without hitting the YC API.

## Regenerating proto code

After editing `proto/compete/v1/compete.proto`:

```bash
buf generate
```
