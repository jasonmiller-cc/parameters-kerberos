FROM golang:1.22-alpine AS builder

WORKDIR /src

# Download dependencies first for layer caching.
COPY go.mod go.sum ./
# The replace directive points to ../parameters-core; copy it alongside.
COPY . .

RUN go build -o /app/parameters-kerberos ./cmd/server

# ---------------------------------------------------------------------------
FROM alpine:3.19 AS runtime

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/parameters-kerberos .

EXPOSE 8080

ENTRYPOINT ["/app/parameters-kerberos"]
