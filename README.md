# 🏗️ Steiger 🏗️

Steiger (dutch for scaffold) is a personal scaffolding for web server projects.

Steiger was initially generated with [Melkeydev/go-blueprint](https://github.com/Melkeydev/go-blueprint) and further modified with my own preferences as well as some
inspiration from [Mat Ryer's](https://bsky.app/profile/matryer.bsky.social) writings on [writing HTTP services in Go](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

live reload the application

```bash
make watch
```

run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```
