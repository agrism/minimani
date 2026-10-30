package game

import (
	"fmt"
	"math"
)

const (
	// Multiple Choice (4 varianti) konfigurācija
	ChoiceBasePoints    = 500
	ChoiceMaxSpeedBonus = 500

	// Exact Input (brīvā ievade) konfigurācija
	ExactMaxBasePoints   = 1000
	ExactBullseyeBonus   = 250
	ExactBullseyeMargin  = 0.05 // 5%
	ExactMaxSpeedBonus   = 150
	ExactSpeedThresholdS = 10
)

type ScoreResult struct {
	GuessValue    float64 `json:"guess_value"`
	ActualPrice   float64 `json:"actual_price"`
	AccuracyPct   float64 `json:"accuracy_pct"`
	BasePoints    int     `json:"base_points"`
	BullseyeBonus int     `json:"bullseye_bonus"`
	SpeedBonus    int     `json:"speed_bonus"`
	TotalPoints   int     `json:"total_points"`
	IsBullseye    bool    `json:"is_bullseye"`
}

type ScoringRuleInfo struct {
	Mode               string   `json:"mode"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	BasePoints         int      `json:"base_points"`
	MaxSpeedBonus      int      `json:"max_speed_bonus"`
	MaxTotalPoints     int      `json:"max_total_points"`
	FormulaExplanation string   `json:"formula_explanation"`
	Examples           []string `json:"examples"`
}

type ScoringRulesResponse struct {
	ActiveMode   string            `json:"active_mode"`
	CurrentRules ScoringRuleInfo   `json:"current_rules"`
	AllModes     []ScoringRuleInfo `json:"all_modes"`
	GeneralTips  []string          `json:"general_tips"`
}

// GetScoringRules returns real-time rules and scoring parameters directly linked to game logic
func GetScoringRules() ScoringRulesResponse {
	choiceRule := ScoringRuleInfo{
		Mode:        "multiple_choice",
		Title:       "4 Cenu Varianti (Multiple Choice)",
		Description: "Katrā raundā tiek doti 4 cenu varianti. Tev pēc sludinājuma foto un parametriem jāizvēlas pareizā ss.lv cena.",
		BasePoints:  ChoiceBasePoints,
		MaxSpeedBonus: ChoiceMaxSpeedBonus,
		MaxTotalPoints: ChoiceBasePoints + ChoiceMaxSpeedBonus,
		FormulaExplanation: fmt.Sprintf(
			"Pareiza atbilde garantē %d bāzes punktus. Papildus vari iegūt līdz pat %d ātruma bonusa punktiem pēc formulas: Ātruma bonuss = (Atlikušais laiks / Kopējais laiks) * %d. Nepareiza atbilde dod 0 punktus.",
			ChoiceBasePoints, ChoiceMaxSpeedBonus, ChoiceMaxSpeedBonus,
		),
		Examples: []string{
			fmt.Sprintf("Atbildēji pareizi uzreiz 1. sekundē: ~%d punkti (%d bāze + ~%d ātrums)", ChoiceBasePoints+ChoiceMaxSpeedBonus-11, ChoiceBasePoints, ChoiceMaxSpeedBonus-11),
			fmt.Sprintf("Atbildēji pareizi raunda vidū (puse laika): ~%d punkti (%d bāze + ~%d ātrums)", ChoiceBasePoints+ChoiceMaxSpeedBonus/2, ChoiceBasePoints, ChoiceMaxSpeedBonus/2),
			fmt.Sprintf("Atbildēji pareizi pēdējā sekundē: %d punkti (%d bāze + 0 ātrums)", ChoiceBasePoints, ChoiceBasePoints),
			"Atbildēji nepareizi: 0 punkti",
		},
	}

	exactRule := ScoringRuleInfo{
		Mode:        "exact",
		Title:       "Precīzā Cenas Ievade (Exact Price)",
		Description: "Spēlētājs pats ar klaviatūru ievada minēto cenu eiro.",
		BasePoints:  ExactMaxBasePoints,
		MaxSpeedBonus: ExactMaxSpeedBonus,
		MaxTotalPoints: ExactMaxBasePoints + ExactBullseyeBonus + ExactMaxSpeedBonus,
		FormulaExplanation: fmt.Sprintf(
			"Punkti aprēķinās pēc simetriskās precizitātes formulas: Precizitāte = (1 - |minējums - faktiskā| / max(minējums, faktiskā)) * 100%%. Bāzes punkti ir līdz %d. Ja kļūda nepārsniedz %.0f%% (Bullseye), saņem papildus +%d punktus. Ja atbildi ātrāk par %ds ar precizitāti virs 70%%, saņem līdz +%d ātruma bonusu.",
			ExactMaxBasePoints, ExactBullseyeMargin*100, ExactBullseyeBonus, ExactSpeedThresholdS, ExactMaxSpeedBonus,
		),
		Examples: []string{
			fmt.Sprintf("100%% Bullseye atbilde pirmajās 2 sekundēs: ~%d punkti", ExactMaxBasePoints+ExactBullseyeBonus+ExactMaxSpeedBonus),
			"80% precīzs minējums: ~800 punkti",
			"50% precīzs minējums: ~500 punkti",
		},
	}

	return ScoringRulesResponse{
		ActiveMode:   "multiple_choice",
		CurrentRules: choiceRule,
		AllModes:     []ScoringRuleInfo{choiceRule, exactRule},
		GeneralTips: []string{
			"Jo ātrāk izvēlies pareizo atbildi, jo vairāk punktu iegūsti — katra sekunde ir svarīga!",
			"Pievērs uzmanību īpašuma parametriem: platība, istabu skaits, atrašanās vieta un ērtības ir galvenie cenu ietekmējošie faktori.",
			"Ja šaubies, salīdzini visus 4 piedāvātos variantus un izslēdz acīmredzami neatbilstošos.",
		},
	}
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

	basePoints := int(math.Round((acc / 100.0) * float64(ExactMaxBasePoints)))

	isBullseye := errRatio <= ExactBullseyeMargin
	bullseyeBonus := 0
	if isBullseye {
		bullseyeBonus = ExactBullseyeBonus
	}

	speedBonus := 0
	thresholdMs := ExactSpeedThresholdS * 1000
	if timeTakenMs < thresholdMs && acc >= 70.0 {
		ratio := float64(thresholdMs-timeTakenMs) / float64(thresholdMs)
		speedBonus = int(math.Round(ratio * float64(ExactMaxSpeedBonus)))
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

	basePoints := ChoiceBasePoints

	// Speed bonus up to ChoiceMaxSpeedBonus points (Kahoot style)
	limitMs := timeLimitSec * 1000
	if limitMs <= 0 {
		limitMs = 45000
	}
	remaining := float64(limitMs - timeTakenMs)
	if remaining < 0 {
		remaining = 0
	}
	speedBonus := int(math.Round((remaining / float64(limitMs)) * float64(ChoiceMaxSpeedBonus)))

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

