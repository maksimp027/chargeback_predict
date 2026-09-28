package engine

import (
	"os"
	"testing"
)

func TestONNXEngine(t *testing.T) {
	dllPath := os.Getenv("ONNX_LIB_PATH")
	if dllPath == "" {
		dllPath = "../../libonnxruntime.so" // Or just skip test
	}
	modelPath := os.Getenv("MODEL_PATH")
	if modelPath == "" {
		modelPath = "../../model_artifacts/chargeback_model.onnx"
	}

	err := InitRuntime(dllPath)
	if err != nil {
		t.Fatalf("Failed to init runtime: %v", err)
	}
	defer CleanupRuntime()

	engine, err := NewONNXEngine(modelPath)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	// 19 dummy features
	dummyFeatures := make([]float32, 19)
	score, err := engine.Predict(dummyFeatures)
	if err != nil {
		t.Fatalf("Prediction failed: %v", err)
	}

	t.Logf("Prediction successful! Score: %f", score)
}
