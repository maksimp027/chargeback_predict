# 🛡️ Chargeback Risk Engine

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.14+-3776AB?style=flat&logo=python)](https://www.python.org/)
[![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![ONNX Runtime](https://img.shields.io/badge/ONNX-Runtime-005CED?style=flat&logo=onnx)](https://onnxruntime.ai/)
[![Redis](https://img.shields.io/badge/Redis-7.0+-DC382D?style=flat&logo=redis)](https://redis.io/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> A high-performance, real-time transaction evaluation microservice designed to detect and prevent chargebacks before they happen.

The **Chargeback Risk Engine** is a robust backend service that evaluates pre-authorization transaction requests via gRPC in microseconds. By combining advanced Machine Learning (XGBoost via ONNX) with real-time streaming velocity metrics (Redis Lua scripts) and business rules, it delivers accurate risk decisions (`APPROVE`, `REVIEW_3DS`, `DECLINE`) along with predictive risk scores.

---

## ✨ Key Features

- **⚡ Blazing Fast gRPC API**: Built for high throughput and ultra-low latency transaction evaluation.
- **🧠 Native ML Inference**: Utilizes C-bindings for ONNX Runtime (`onnxruntime_go`) to execute a pre-trained XGBoost model (`chargeback_model.onnx`) directly in Go.
- **🔄 Real-Time Velocity Tracking**: Leverages custom Redis Lua scripts for atomic, sliding-window feature extraction (e.g., transaction counts and amounts over 5m, 15m, 60m windows).
- **🛡️ Hybrid Rule Engine**: Combines the ML probability score with hard business constraints (CVV mismatches, 3D-Secure state, velocity limits) to produce actionable decisions.
- **📊 Observability Ready**: Easily integrates with Prometheus & Grafana (configured via Docker Compose).

---

## 🏗️ Architecture & Data Flow

1. **Ingest**: Client sends a `RiskEvaluationRequest` to the gRPC server.
2. **Velocity Aggregation**: The velocity client runs a Lua script against Redis to retrieve sliding window metrics.
3. **Feature Enrichment**: Combines request metadata and velocity metrics into a 19-element feature vector.
4. **ML Inference**: The ONNX engine predicts the chargeback probability.
5. **Decisioning**: The Rule Engine applies business logic thresholds over the ML score to output the final decision and reason codes.

---

## 🛠️ Tech Stack

- **Core Service:** Go 1.22+
- **Machine Learning:** Python (XGBoost, scikit-learn), ONNX
- **Inference Engine:** ONNX Runtime C Shared Library
- **Data Store / Caching:** Redis 7+
- **Communication:** gRPC / Protocol Buffers
- **Infrastructure:** Docker & Docker Compose

---

## 📂 Project Structure

```text
chargeback-risk-engine/
├── api/proto/             # Protocol Buffer definitions & generated gRPC code
├── cmd/server/            # Application entrypoint
├── deploy/                # Dockerfile, docker-compose.yml, and Prometheus config
├── internal/
│   ├── config/            # Environment configuration loader
│   ├── domain/            # Domain models and schema definitions
│   ├── engine/            # ONNX Runtime integration wrapper
│   ├── features/          # Feature extraction & Redis velocity client
│   ├── rules/             # Decision threshold engine
│   └── service/           # gRPC service implementation
├── ml/                    # Python training, feature extraction scripts, and ONNX export
├── model_artifacts/       # Trained ONNX model and feature schema JSON
└── Makefile               # Task automation (build, test, docker, etc.)
```

---

## 🚀 Getting Started

### Prerequisites
- Go 1.22 or higher
- Python 3.14+ (for retraining models)
- Docker & Docker Compose
- ONNX Runtime C shared library

### Quick Start (Docker)

To spin up the service along with Redis and Prometheus:

```bash
make docker-up
```
*The gRPC server will start listening on port `50051`.*

### Local Development Setup

1. **Set up the ONNX Runtime library:**
   Ensure the ONNX Runtime library path is set.
   ```bash
   # Linux Example (if installed system-wide)
   export ONNX_LIB_PATH=/usr/lib/libonnxruntime.so
   
   # Or using the local Python venv (if installed via pip in this project)
   export ONNX_LIB_PATH=./libonnxruntime.so
   ```

2. **Run Unit Tests:**
   ```bash
   make test
   ```

3. **Build the Binary:**
   ```bash
   make build
   ```

4. **Run the Server:**
   ```bash
   ./bin/chargeback-risk-engine
   ```

---

## ⚙️ Configuration

The service is highly configurable via environment variables:

| Variable | Default Value | Description |
|---|---|---|
| `SERVER_ADDR` | `:50051` | gRPC server bind address |
| `REDIS_ADDR` | `localhost:6379` | Redis server address for velocity metrics |
| `MODEL_PATH` | `model_artifacts/chargeback_model.onnx` | Path to the ONNX model file |
| `SCHEMA_PATH` | `model_artifacts/features_schema.json` | Path to the feature schema definition |
| `ONNX_LIB_PATH` | `/usr/lib/libonnxruntime.so` | Path to ONNX Runtime shared library |

---

## 🧑‍💻 ML Pipeline (Python)

The `ml/` directory contains everything needed to train the XGBoost model.
To set up the Python environment and retrain:

```bash
# Create venv and activate
python3 -m venv venv
source venv/bin/activate

# Install dependencies
pip install -r ml/requirements.txt

# Run training
python ml/train.py

# Export to ONNX
python ml/export_onnx.py
```

---

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
