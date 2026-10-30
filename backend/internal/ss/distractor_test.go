package ss

import (
	"testing"
)

func TestGeneratePriceOptions(t *testing.T) {
	actual := 12500.0
	options := GeneratePriceOptions(actual)

	if len(options) != 4 {
		t.Fatalf("Expected 4 options, got %d", len(options))
	}

	foundActual := false
	seen := make(map[float64]bool)
	for _, opt := range options {
		if opt == actual {
			foundActual = true
		}
		if seen[opt] {
			t.Errorf("Duplicate option found: %.2f", opt)
		}
		seen[opt] = true
	}

	if !foundActual {
		t.Errorf("Actual price %.2f was not found in generated options", actual)
	}
}

