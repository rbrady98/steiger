FROM golang:1.24-bookworm AS base

# Development stage
# Meant to be used with compose so doesnt include the files
FROM base AS development

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN go mod download

CMD ["air"]

# Builder stage
FROM base AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

# Production stage
FROM scratch AS production

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=builder /build/build/app ./

CMD ["./app"]
