# Chargeback Risk Engine

High-performance real-time chargeback risk evaluation service written in Go, powered by an ONNX-exported XGBoost model and Redis velocity tracking.

## Overview

The service evaluates incoming pre-authorization transaction requests via gRPC and returns a risk score, decision (`APPROVE`, `REVIEW_3DS`, `DECLINE`), execution latency in microseconds, and applicable reason codes.

### Key Components

- **gRPC API**: High-throughput transaction evaluation endpoint (`risk.v1.RiskService`).
- **ONNX Inference Engine**: Native C-bindings (`onnxruntime_go`) executing pre-trained XGBoost model (`chargeback_model.onnx`).
- **Real-Time Velocity Tracking**: Atomic sliding-window feature extraction in Redis using custom Lua script (`card_tx_cnt_5m`, `15m`, `60m`, `card_sum_cents_60m`).
- **Feature Builder**: Dynamic feature mapping and normalization matching the offline ML pipeline.
- **Rule Engine**: Combines model probabilities with rule-based thresholds and risk factors (CVV verification, 3D-Secure, velocity limits).

## Architecture & Data Flow

1. Client sends `RiskEvaluationRequest` to gRPC server.
2. Velocity client executes Lua script against Redis to retrieve sliding window metrics.
3. Feature enricher combines request metadata and velocity metrics into a 19-element float32 array.
4. ONNX engine executes model inference and returns chargeback probability.
5. Rule engine evaluates ML score alongside business rules to produce the final decision and reason codes.

## Requirements

- Go 1.22+
- Docker & Docker Compose
- ONNX Runtime C shared library (`libonnxruntime.so` on Linux or `onnxruntime.dll` on Windows)
- Redis 7+

## Project Structure

```
chargeback-risk-engine/
├── api/proto/             # Protocol Buffer definitions and generated gRPC code
├── cmd/server/            # Application entrypoint
├── deploy/                # Dockerfile, docker-compose.yml, and Prometheus config
├── internal/
│   ├── config/            # Environment configuration loader
│   ├── domain/            # Domain models and schema definitions
│   ├── engine/            # ONNX Runtime integration wrapper
│   ├── features/          # Feature extraction and Redis velocity client
│   ├── rules/             # Decision threshold engine
│   └── service/           # gRPC service implementation
├── ml/                    # Python training, feature script, and ONNX export
├── model_artifacts/       # Trained ONNX model and feature schema JSON
├── Makefile               # Task automation
└── README.md
```

## Quick Start

### Running with Docker Compose

To spin up the service along with Redis and Prometheus:

```bash
make docker-up
```

Or directly via docker compose:

```bash
docker compose -f deploy/docker-compose.yml up --build -d
```

The gRPC server will start listening on port `50051`.

### Local Development

1. Ensure ONNX Runtime library path is set:
   - Linux: `export ONNX_LIB_PATH=/usr/lib/libonnxruntime.so`
   - Windows: `set ONNX_LIB_PATH=C:\path\to\onnxruntime.dll`

2. Run unit tests:

```bash
make test
```

3. Build the binary locally:

```bash
make build
```

4. Run the server:

```bash
./bin/chargeback-risk-engine
```

## Configuration

The service can be configured using environment variables:

| Variable | Default Value | Description |
|---|---|---|
| `SERVER_ADDR` | `:50051` | gRPC server bind address |
| `REDIS_ADDR` | `localhost:6379` | Redis server address for velocity metrics |
| `MODEL_PATH` | `model_artifacts/chargeback_model.onnx` | Path to ONNX model file |
| `SCHEMA_PATH` | `model_artifacts/features_schema.json` | Path to feature schema definition |
| `ONNX_LIB_PATH` | `/usr/lib/libonnxruntime.so` | Path to ONNX Runtime shared library |

## License

MIT
