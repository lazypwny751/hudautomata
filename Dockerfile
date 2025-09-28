## Simple multi-stage build (Debian runtime as requested)
FROM golang:1.23.0 AS builder
WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Source
COPY . .

# Build (static)
ENV CGO_ENABLED=0
RUN go build -o hudautomata "src/main.go"

FROM debian:12-slim
WORKDIR /app

# Copy binary only
COPY --from=builder /app/hudautomata /usr/local/bin/hudautomata

EXPOSE 8080

ENV GIN_MODE=release

ENTRYPOINT ["hudautomata"]
CMD ["--host", "0.0.0.0", "--port", "8080"]
