package ss

import (
	"encoding/json"
	"math"
	"math/rand"
)

// GeneratePriceOptions generates 4 plausible price options (1 correct, 3 distractors)
// formatted as realistic Latvian classified numbers.
func GeneratePriceOptions(actualPrice float64) []float64 {
	if actualPrice <= 0 {
		return []float64{100, 250, 500, 1000}
	}

	distractorRatios := []float64{
		0.45 + rand.Float64()*0.20, // 0.45 - 0.65x
		1.35 + rand.Float64()*0.35, // 1.35 - 1.70x
		2.10 + rand.Float64()*0.70, // 2.10 - 2.80x
	}

	// Sometimes have one slightly lower and two higher or vice versa
	if rand.Float64() < 0.3 {
		distractorRatios = []float64{
			0.25 + rand.Float64()*0.15, // very cheap 0.25 - 0.40x
			0.65 + rand.Float64()*0.15, // 0.65 - 0.80x
			1.50 + rand.Float64()*0.40, // 1.50 - 1.90x
		}
	}

	options := make([]float64, 0, 4)
	options = append(options, actualPrice)

	seen := make(map[float64]bool)
	seen[actualPrice] = true

	for _, ratio := range distractorRatios {
		fake := roundToCleanPrice(actualPrice * ratio)
		if fake <= 0 || seen[fake] {
			fake = roundToCleanPrice(actualPrice * (ratio + 0.25))
		}
		seen[fake] = true
		options = append(options, fake)
	}

	// Shuffle the 4 options
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return options
}

// GenerateOptionsJSON returns the options as JSON bytes
func GenerateOptionsJSON(actualPrice float64) []byte {
	opts := GeneratePriceOptions(actualPrice)
	data, _ := json.Marshal(opts)
	return data
}

// roundToCleanPrice rounds prices to realistic marketing numbers like 4 900, 12 500, etc.
func roundToCleanPrice(val float64) float64 {
	if val < 50 {
		return math.Round(val)
	}
	if val < 500 {
		// Round to nearest 5 or 10
		return math.Round(val/10) * 10
	}
	if val < 2000 {
		// Round to nearest 50
		return math.Round(val/50) * 50
	}
	if val < 10000 {
		// Round to nearest 100 or 250
		return math.Round(val/100) * 100
	}
	if val < 50000 {
		// Round to nearest 500
		return math.Round(val/500) * 500
	}
	// Round to nearest 1000
	return math.Round(val/1000) * 1000
}

