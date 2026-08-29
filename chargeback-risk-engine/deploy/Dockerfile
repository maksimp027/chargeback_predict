# Multi-stage Dockerfile for Chargeback Risk Engine

# Stage 1: Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /src

# Download and install ONNX Runtime library for Linux x86_64
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    tar \
    build-essential \
 && rm -rf /var/lib/apt/lists/*

ENV ONNX_VERSION=1.20.1
RUN curl -sSL "https://github.com/microsoft/onnxruntime/releases/download/v${ONNX_VERSION}/onnxruntime-linux-x64-${ONNX_VERSION}.tgz" \
    | tar -xz -C /tmp \
 && cp /tmp/onnxruntime-linux-x64-${ONNX_VERSION}/lib/libonnxruntime.so* /usr/lib/ \
 && ldconfig

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/server ./cmd/server

# Stage 2: Runtime stage
FROM debian:bookworm-slim AS runner

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

# Copy ONNX Runtime shared library from builder stage
COPY --from=builder /usr/lib/libonnxruntime.so* /usr/lib/
RUN ldconfig

# Copy compiled binary and model artifacts
COPY --from=builder /app/server /app/server
COPY model_artifacts /app/model_artifacts

# Set default environment variables
ENV SERVER_ADDR=":50051" \
    REDIS_ADDR="redis:6379" \
    MODEL_PATH="/app/model_artifacts/chargeback_model.onnx" \
    SCHEMA_PATH="/app/model_artifacts/features_schema.json" \
    ONNX_LIB_PATH="/usr/lib/libonnxruntime.so"

EXPOSE 50051

ENTRYPOINT ["/app/server"]
