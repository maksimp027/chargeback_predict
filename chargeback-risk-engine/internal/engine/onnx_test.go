package engine

import (
	"testing"
)

func TestONNXEngine(t *testing.T) {
	dllPath := `D:\golandproject\venv\Lib\site-packages\onnxruntime\capi\onnxruntime.dll`
	modelPath := `D:\golandproject\chargeback-risk-engine\model_artifacts\chargeback_model.onnx`

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
