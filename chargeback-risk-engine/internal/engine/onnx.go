package engine

import (
	"fmt"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
)

type ONNXEngine struct {
	session      *ort.AdvancedSession
	inputTensor  *ort.Tensor[float32]
	outputTensor *ort.Tensor[float32]
	mu           sync.Mutex
}

func InitRuntime(sharedLibPath string) error {
	ort.SetSharedLibraryPath(sharedLibPath)
	return ort.InitializeEnvironment()
}

func CleanupRuntime() {
	ort.DestroyEnvironment()
}

// NewONNXEngine creates an ONNX engine and prepares the session and pre-allocated tensors
func NewONNXEngine(modelPath string) (*ONNXEngine, error) {
	inputNames := []string{"float_input"}
	outputNames := []string{"probabilities"}

	// Expected data shape: 1 transaction, 19 features
	inputShape := ort.NewShape(1, 19)
	// Output shape: 1 transaction, 2 classes [P(Normal), P(Chargeback)]
	outputShape := ort.NewShape(1, 2)

	// Pre-allocate input tensor
	inputData := make([]float32, 19)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	if err != nil {
		return nil, fmt.Errorf("failed to create input tensor: %w", err)
	}

	// Pre-allocate output tensor
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		inputTensor.Destroy()
		return nil, fmt.Errorf("failed to create output tensor: %w", err)
	}

	session, err := ort.NewAdvancedSession(
		modelPath,
		inputNames,
		outputNames,
		[]ort.Value{inputTensor},
		[]ort.Value{outputTensor},
		nil,
	)
	if err != nil {
		inputTensor.Destroy()
		outputTensor.Destroy()
		return nil, fmt.Errorf("failed to create ONNX session: %w", err)
	}

	return &ONNXEngine{
		session:      session,
		inputTensor:  inputTensor,
		outputTensor: outputTensor,
	}, nil
}

// Predict performs risk calculation (returns float32 from 0.0 to 1.0)
func (e *ONNXEngine) Predict(features []float32) (float32, error) {
	if len(features) != 19 {
		return 0.0, fmt.Errorf("invalid feature length: expected 19, got %d", len(features))
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Copy new features into pre-allocated input tensor memory
	copy(e.inputTensor.GetData(), features)

	// Run inference
	err := e.session.Run()
	if err != nil {
		return 0.0, fmt.Errorf("ONNX inference execution failed: %w", err)
	}

	// Read the result from pre-allocated output tensor memory
	results := e.outputTensor.GetData()
	if len(results) < 2 {
		return 0.0, fmt.Errorf("unexpected output shape from model")
	}

	// results[0] = P(Normal), results[1] = P(Chargeback)
	return results[1], nil
}

func (e *ONNXEngine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var err error
	if e.session != nil {
		err = e.session.Destroy()
	}
	if e.inputTensor != nil {
		e.inputTensor.Destroy()
	}
	if e.outputTensor != nil {
		e.outputTensor.Destroy()
	}
	return err
}
