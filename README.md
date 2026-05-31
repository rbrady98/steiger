# 🏗️ Steiger 🏗️

Steiger (dutch for scaffold) is a personal scaffolding for web server projects.

Steiger draws inspiration from [Mat Ryer's](https://bsky.app/profile/matryer.bsky.social) writings on [writing HTTP services in Go](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years)
and the [DDD-lite in Go](https://threedots.tech/post/ddd-lite-in-go-introduction/) approach from Three Dots Labs.

## Motivation

- Sensible, out-of-the-box defaults for Go HTTP APIs.
- A deliberate reaction against overengineering and the overuse of abstractions.
- Clear separation of concerns utilizing a "DDD-lite" approach, resulting in an intuitive, easily navigable package structure.
- Leveraging code-generation tools like sqlc (and eventually oapi-codegen) to keep boilerplate to an absolute minimum.
## Getting Started

Common tasks are exposed through the `Makefile`:

| Command | Description |
| --- | --- |
| `make build` | Build the binary to `build/main` |
| `make run` | Run the API |
| `make dev` | Live-reload via `air` |
| `make test` | Run the test suite |
| `make lint` | Run `golangci-lint` |
| `make sqlcgen` | Regenerate sqlc query bindings |
| `make clean` | Remove build artifacts |

## OpenTelemetry
If you want to test the otel setup you can locally run the "LGTM" stack. Navigate to localhost:3000 to view your telemetry
```bash
docker run --rm --network="host" grafana/otel-lgtm:latest
```
