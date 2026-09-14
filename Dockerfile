# syntax=docker/dockerfile:1
FROM golang:1.27-alpine AS builder

WORKDIR /src

# Download dependencies first for layer caching.
COPY go.mod go.sum ./
# parameters-core is a local replace living in a sibling directory outside
# this build context; Docker forbids COPY-ing paths outside the primary
# context, so it's supplied as a named build context instead — see
# --build-context in the Makefile.
COPY --from=parameters-core . /parameters-core
COPY . .

RUN go build -o /app/parameters-kerberos ./cmd/server

# ---------------------------------------------------------------------------
FROM alpine:3.19 AS runtime

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/parameters-kerberos .

EXPOSE 8080

ENTRYPOINT ["/app/parameters-kerberos"]
