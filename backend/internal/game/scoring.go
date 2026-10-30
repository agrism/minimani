package game

import (
	"math"
)

type ScoreResult struct {
	GuessValue   float64 `json:"guess_value"`
	ActualPrice  float64 `json:"actual_price"`
	AccuracyPct  float64 `json:"accuracy_pct"`
	BasePoints   int     `json:"base_points"`
	BullseyeBonus int    `json:"bullseye_bonus"`
	SpeedBonus   int     `json:"speed_bonus"`
	TotalPoints  int     `json:"total_points"`
	IsBullseye   bool    `json:"is_bullseye"`
}

// CalculateExactScore computes points for exact price input
func CalculateExactScore(guess, actual float64, timeTakenMs int, timeLimitSec int) ScoreResult {
	if guess <= 0 || actual <= 0 {
		return ScoreResult{
			GuessValue:  guess,
			ActualPrice: actual,
			AccuracyPct: 0,
			TotalPoints: 0,
		}
	}

	maxVal := math.Max(guess, actual)
	diff := math.Abs(guess - actual)
	errRatio := diff / maxVal

	acc := math.Max(0, (1.0-errRatio)*100.0)
	acc = math.Round(acc*100) / 100 // round to 2 decimals

	basePoints := int(math.Round(acc * 10)) // 0 to 1000 points

	isBullseye := errRatio <= 0.05 // +-5% margin
	bullseyeBonus := 0
	if isBullseye {
		bullseyeBonus = 250
	}

	speedBonus := 0
	// If answered quickly (under 10s) and accuracy is >= 70%
	if timeTakenMs < 10000 && acc >= 70.0 {
		ratio := float64(10000-timeTakenMs) / 10000.0
		speedBonus = int(math.Round(ratio * 150))
	}

	total := basePoints + bullseyeBonus + speedBonus

	return ScoreResult{
		GuessValue:    guess,
		ActualPrice:   actual,
		AccuracyPct:   acc,
		BasePoints:    basePoints,
		BullseyeBonus: bullseyeBonus,
		SpeedBonus:    speedBonus,
		TotalPoints:   total,
		IsBullseye:    isBullseye,
	}
}

// CalculateChoiceScore computes points for 4-option multiple choice mode
func CalculateChoiceScore(guess, actual float64, timeTakenMs int, timeLimitSec int) ScoreResult {
	isCorrect := math.Abs(guess-actual) < 0.01

	if !isCorrect {
		return ScoreResult{
			GuessValue:  guess,
			ActualPrice: actual,
			AccuracyPct: 0,
			TotalPoints: 0,
		}
	}

	basePoints := 500

	// Speed bonus up to 500 points (Kahoot style)
	limitMs := timeLimitSec * 1000
	if limitMs <= 0 {
		limitMs = 30000
	}
	remaining := float64(limitMs - timeTakenMs)
	if remaining < 0 {
		remaining = 0
	}
	speedBonus := int(math.Round((remaining / float64(limitMs)) * 500))

	return ScoreResult{
		GuessValue:    guess,
		ActualPrice:   actual,
		AccuracyPct:   100.0,
		BasePoints:    basePoints,
		BullseyeBonus: 0,
		SpeedBonus:    speedBonus,
		TotalPoints:   basePoints + speedBonus,
		IsBullseye:    true,
	}
}

