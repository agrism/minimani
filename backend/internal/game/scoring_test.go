package game

import (
	"testing"
)

func TestCalculateExactScore(t *testing.T) {
	// 1. Exact match (100% accuracy, Bullseye bonus, speed bonus)
	res := CalculateExactScore(10000, 10000, 3000, 45)
	if res.AccuracyPct != 100.0 {
		t.Errorf("Expected 100.0 accuracy, got %.2f", res.AccuracyPct)
	}
	if !res.IsBullseye {
		t.Errorf("Expected IsBullseye to be true")
	}
	if res.BullseyeBonus != 250 {
		t.Errorf("Expected 250 bullseye bonus, got %d", res.BullseyeBonus)
	}
	if res.TotalPoints <= 1250 {
		t.Errorf("Expected total points > 1250 with speed bonus, got %d", res.TotalPoints)
	}

	// 2. Off by 20% (8000 vs 10000)
	res2 := CalculateExactScore(8000, 10000, 15000, 45)
	if res2.AccuracyPct != 80.0 {
		t.Errorf("Expected 80.0 accuracy, got %.2f", res2.AccuracyPct)
	}
	if res2.IsBullseye {
		t.Errorf("Expected IsBullseye to be false")
	}
	if res2.BasePoints != 800 {
		t.Errorf("Expected 800 base points, got %d", res2.BasePoints)
	}
}

func TestCalculateChoiceScore(t *testing.T) {
	// Correct choice
	res := CalculateChoiceScore(4500, 4500, 5000, 30)
	if res.TotalPoints <= 500 {
		t.Errorf("Expected total points > 500 with speed bonus, got %d", res.TotalPoints)
	}

	// Incorrect choice
	resFail := CalculateChoiceScore(2000, 4500, 2000, 30)
	if resFail.TotalPoints != 0 {
		t.Errorf("Expected 0 points for incorrect choice, got %d", resFail.TotalPoints)
	}
}

