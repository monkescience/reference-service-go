# AGENTS

Ultra-brief notes so automation and bots can build, run, and touch this service safely.

## Quick facts
- Language/runtime: Go (see `go.mod`). Preferred way to build is via Docker to avoid local Go version drift.
- HTTP port: 8080
- Base API path: `/v1`
- Required env: `VERSION` (process will exit if missing). No other config.
- Health/metrics: `/status/live`, `/status/ready`, `/status/metrics` (Prometheus)
- OpenAPI source: `openapi/v1_order_api.yaml`
- Codegen target package: `internal/adapters/http/order` (files `models.gen.go`, `server.gen.go`)

## Build and run
- Docker (recommended):
```sh
# from repo root
export VERSION=0.0.0
docker build -t reference-service-go:dev --build-arg VERSION="$VERSION" .
docker run --rm -p 8080:8080 -e VERSION="$VERSION" reference-service-go:dev
```
- Local (if you have a compatible Go toolchain):
```sh
export VERSION=0.0.0
go run ./cmd/main.go
```

## API surface
- Orders:
  - GET `/v1/orders?customer_id=<uuid>&limit=<int>&offset=<int>`
  - POST `/v1/orders` (JSON body per OpenAPI)
  - GET `/v1/orders/{order_id}`
- Health & metrics:
  - GET `/status/live`, `/status/ready`, `/status/metrics`
- Schema is defined in `openapi/v1_order_api.yaml`. Generated types/routers live under `internal/adapters/http/order`.

## Regenerate HTTP code from OpenAPI
- The repo embeds `go:generate` directives in `internal/tools.go`.
```sh
# Regenerate models and server stubs
go generate ./internal
```
This uses `github.com/oapi-codegen/oapi-codegen/v2` and configs in `openapi/oapi-codegen.*.yaml`.

## Observability
- Prometheus metric: `app_http_server_request_duration_seconds` (histogram) with labels `method`, `route`, `code`.
- Metrics are exposed at `/status/metrics`.

## Minimal architecture
- Onion architecture: domain (`internal/domain`), use cases (`internal/usecase`), ports (`internal/ports`), adapters (`internal/adapters`).
- Default repository is in-memory: `internal/adapters/repository/memory` implementing `internal/ports/order.Repository`.

## Gotchas for agents
- Always set `VERSION` when starting the binary (the app will `log.Fatal` if it is empty).
- Prefer Docker builds to sidestep local Go version mismatches.
- If you touch OpenAPI, re-run codegen and commit the updated `*.gen.go` files.

## Smoke test
```sh
# Place an order
curl -sS -X POST http://localhost:8080/v1/orders \
  -H 'content-type: application/json' \
  -d '{"customer_id":"00000000-0000-0000-0000-000000000000","items":[{"name":"widget"}]}' | jq .

# List orders
curl -sS 'http://localhost:8080/v1/orders?limit=10' | jq .

# Health
curl -sS http://localhost:8080/status/live | jq .
```
