package game

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	internaldb "minimani/backend/internal/db"
	db "minimani/backend/internal/db/sqlc"
	"minimani/backend/internal/ss"
)

type RoomStatus string

const (
	StatusLobby       RoomStatus = "lobby"
	StatusStarting    RoomStatus = "starting"
	StatusRoundActive RoomStatus = "round_active"
	StatusRoundReveal RoomStatus = "round_reveal"
	StatusGameOver    RoomStatus = "game_over"
)

type Player struct {
	ID        string `json:"id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsReady   bool   `json:"is_ready"`
	Score     int    `json:"score"`
	HasVoted  bool   `json:"has_voted"`
	HasGuess  bool   `json:"has_guess"`
}

type VoteSummary struct {
	Options map[string]int     `json:"options"` // e.g. "10": 3
	Percents map[string]float64 `json:"percents"` // e.g. "10": 60.0
	Winning string             `json:"winning"`
}

type RoundListing struct {
	ID          string    `json:"id"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImagePath   string    `json:"image_path"`
	Region      string    `json:"region"`
	Options     []float64 `json:"options,omitempty"` // 4 choices for multiple_choice mode
	ActualPrice float64   `json:"actual_price,omitempty"` // ONLY present in round_reveal
}

type PlayerGuessResult struct {
	PlayerID    string  `json:"player_id"`
	Nickname    string  `json:"nickname"`
	Avatar      string  `json:"avatar"`
	GuessValue  float64 `json:"guess_value"`
	AccuracyPct float64 `json:"accuracy_pct"`
	ScorePoints int     `json:"score_points"`
	IsBullseye  bool    `json:"is_bullseye"`
	TimeTakenMs int     `json:"time_taken_ms"`
}

type Room struct {
	mu sync.RWMutex

	ID            string                `json:"id"`
	Code          string                `json:"code"`
	HostID        string                `json:"host_id"`
	IsPublic      bool                  `json:"is_public"`
	Category      string                `json:"category"`
	GameMode      string                `json:"game_mode"`      // "exact", "multiple_choice", "mix"
	QuestionCount int                   `json:"question_count"` // 5, 10, 15, 20
	TimeLimitSec  int                   `json:"time_limit_sec"`
	MaxPlayers    int                   `json:"max_players"`
	Status        RoomStatus            `json:"status"`

	Players       map[string]*Player    `json:"players"`
	VotesCount    map[string]int        `json:"votes_count"` // player_id -> 5, 10, 15, 20
	VotesMode     map[string]string     `json:"votes_mode"`  // player_id -> exact, multiple_choice
	VotesCategory map[string]string     `json:"votes_category"`

	CurrentRound   int                  `json:"current_round"`
	RoundStartedAt time.Time            `json:"round_started_at"`
	RoundGuesses   map[string]PlayerGuessResult `json:"round_guesses"`

	Listings       []RoundListing       `json:"-"`
	CurrentListing *RoundListing        `json:"current_listing,omitempty"`

	BroadcastFunc  func(msg []byte)     `json:"-"`
	TimerCancel    context.CancelFunc   `json:"-"`
	Queries        *db.Queries          `json:"-"`
}

func NewRoom(code, hostID string, isPublic bool, queries *db.Queries, broadcast func([]byte)) *Room {
	return &Room{
		Code:          code,
		HostID:        hostID,
		IsPublic:      isPublic,
		Category:      "all",
		GameMode:      "exact",
		QuestionCount: 10,
		TimeLimitSec:  45,
		MaxPlayers:    8,
		Status:        StatusLobby,
		Players:       make(map[string]*Player),
		VotesCount:    make(map[string]int),
		VotesMode:     make(map[string]string),
		VotesCategory: make(map[string]string),
		RoundGuesses:  make(map[string]PlayerGuessResult),
		BroadcastFunc: broadcast,
		Queries:       queries,
	}
}

// AddPlayer adds a player to the room
func (r *Room) AddPlayer(id, nickname, avatar string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if avatar == "" {
		avatar = "🚗"
	}
	if nickname == "" {
		nickname = fmt.Sprintf("Viesis #%d", rand.Intn(900)+100)
	}

	r.Players[id] = &Player{
		ID:       id,
		Nickname: nickname,
		Avatar:   avatar,
		IsReady:  false,
		Score:    0,
	}
}

// RemovePlayer removes a player from the room
func (r *Room) RemovePlayer(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.Players, id)
	delete(r.VotesCount, id)
	delete(r.VotesMode, id)
	delete(r.VotesCategory, id)
}

// Vote updates player votes and computes consensus
func (r *Room) Vote(playerID string, questionCount int, gameMode string, category string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Status != StatusLobby {
		return
	}

	if questionCount == 5 || questionCount == 10 || questionCount == 15 || questionCount == 20 {
		r.VotesCount[playerID] = questionCount
	}
	if gameMode == "exact" || gameMode == "multiple_choice" || gameMode == "mix" {
		r.VotesMode[playerID] = gameMode
	}
	if category != "" {
		r.VotesCategory[playerID] = category
	}

	totalPlayers := len(r.Players)
	if totalPlayers == 0 {
		return
	}

	// Calculate consensus for QuestionCount
	countTallies := make(map[int]int)
	for _, v := range r.VotesCount {
		countTallies[v]++
		if float64(countTallies[v])/float64(totalPlayers) > 0.5 {
			r.QuestionCount = v
		}
	}

	// Calculate consensus for GameMode
	modeTallies := make(map[string]int)
	for _, v := range r.VotesMode {
		modeTallies[v]++
		if float64(modeTallies[v])/float64(totalPlayers) > 0.5 {
			r.GameMode = v
		}
	}

	// Calculate consensus for Category
	catTallies := make(map[string]int)
	for _, v := range r.VotesCategory {
		catTallies[v]++
		if float64(catTallies[v])/float64(totalPlayers) > 0.5 {
			r.Category = v
		}
	}
}

// SetReady sets player ready status
func (r *Room) SetReady(playerID string, ready bool) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if p, exists := r.Players[playerID]; exists {
		p.IsReady = ready
	}

	// Check if all players are ready and at least 2 players (or 1 in dev/testing)
	if len(r.Players) >= 1 {
		allReady := true
		for _, p := range r.Players {
			if !p.IsReady {
				allReady = false
				break
			}
		}
		return allReady
	}
	return false
}

// StartGame loads listings from database and starts round 1
func (r *Room) StartGame(ctx context.Context) error {
	r.mu.Lock()

	// Load random listings
	dbListings, err := r.Queries.GetRandomListings(ctx, db.GetRandomListingsParams{
		Category: r.Category,
		Limit:    int32(r.QuestionCount),
	})
	if err != nil {
		r.mu.Unlock()
		return fmt.Errorf("failed to fetch listings: %w", err)
	}

	if len(dbListings) == 0 {
		r.mu.Unlock()
		return fmt.Errorf("nav pieejamu sludinājumu šajā kategorijā")
	}

	r.Listings = make([]RoundListing, 0, len(dbListings))
	for _, l := range dbListings {
		var opts []float64
		if len(l.Options) > 0 {
			_ = json.Unmarshal(l.Options, &opts)
		}
		actual := internaldb.FloatFromNumeric(l.ActualPrice)
		if len(opts) == 0 {
			opts = ss.GeneratePriceOptions(actual)
		}

		r.Listings = append(r.Listings, RoundListing{
			ID:          l.ID.String(),
			Category:    l.Category,
			Title:       l.Title,
			Description: l.Description.String,
			ImagePath:   l.ImagePath,
			Region:      l.Region.String,
			Options:     opts,
			ActualPrice: actual,
		})
	}

	r.QuestionCount = len(r.Listings)
	r.CurrentRound = 0
	r.Status = StatusStarting
	r.mu.Unlock()

	r.BroadcastState()

	// 3-second countdown before round 1
	go func() {
		time.Sleep(3 * time.Second)
		r.NextRound()
	}()

	return nil
}

// NextRound starts the next question round
func (r *Room) NextRound() {
	r.mu.Lock()

	if r.CurrentRound >= len(r.Listings) {
		r.Status = StatusGameOver
		r.mu.Unlock()
		r.BroadcastState()
		return
	}

	r.CurrentRound++
	listing := r.Listings[r.CurrentRound-1]

	// Create safe client listing WITHOUT actual_price
	clientListing := RoundListing{
		ID:          listing.ID,
		Category:    listing.Category,
		Title:       listing.Title,
		Description: listing.Description,
		ImagePath:   listing.ImagePath,
		Region:      listing.Region,
	}

	// Determine round mode (if mix, alternate)
	roundMode := r.GameMode
	if roundMode == "mix" {
		if r.CurrentRound%2 == 0 {
			roundMode = "multiple_choice"
		} else {
			roundMode = "exact"
		}
	}

	if roundMode == "multiple_choice" {
		clientListing.Options = listing.Options
	}

	r.CurrentListing = &clientListing
	r.Status = StatusRoundActive
	r.RoundStartedAt = time.Now()
	r.RoundGuesses = make(map[string]PlayerGuessResult)

	for _, p := range r.Players {
		p.HasGuess = false
	}

	ctx, cancel := context.WithCancel(context.Background())
	r.TimerCancel = cancel
	r.mu.Unlock()

	r.BroadcastState()

	// Launch round countdown timer
	go func(roundNum int) {
		select {
		case <-time.After(time.Duration(r.TimeLimitSec) * time.Second):
			r.FinishRound(roundNum)
		case <-ctx.Done():
			// Round ended early because all players guessed
		}
	}(r.CurrentRound)
}

// SubmitGuess records a player's price guess
func (r *Room) SubmitGuess(playerID string, guessValue float64) bool {
	r.mu.Lock()

	if r.Status != StatusRoundActive {
		r.mu.Unlock()
		return false
	}

	player, exists := r.Players[playerID]
	if !exists || player.HasGuess {
		r.mu.Unlock()
		return false
	}

	timeTaken := int(time.Since(r.RoundStartedAt).Milliseconds())
	listing := r.Listings[r.CurrentRound-1]

	var scoreRes ScoreResult
	if len(r.CurrentListing.Options) > 0 {
		scoreRes = CalculateChoiceScore(guessValue, listing.ActualPrice, timeTaken, r.TimeLimitSec)
	} else {
		scoreRes = CalculateExactScore(guessValue, listing.ActualPrice, timeTaken, r.TimeLimitSec)
	}

	player.HasGuess = true
	player.Score += scoreRes.TotalPoints

	r.RoundGuesses[playerID] = PlayerGuessResult{
		PlayerID:    playerID,
		Nickname:    player.Nickname,
		Avatar:      player.Avatar,
		GuessValue:  guessValue,
		AccuracyPct: scoreRes.AccuracyPct,
		ScorePoints: scoreRes.TotalPoints,
		IsBullseye:  scoreRes.IsBullseye,
		TimeTakenMs: timeTaken,
	}

	// Check if all connected players have guessed
	allGuessed := true
	for _, p := range r.Players {
		if !p.HasGuess {
			allGuessed = false
			break
		}
	}

	roundNum := r.CurrentRound
	cancel := r.TimerCancel
	r.mu.Unlock()

	if allGuessed {
		if cancel != nil {
			cancel()
		}
		go r.FinishRound(roundNum)
	} else {
		r.BroadcastState()
	}

	return true
}

// FinishRound reveals actual price, results and updates scores
func (r *Room) FinishRound(roundNum int) {
	r.mu.Lock()

	if r.Status != StatusRoundActive || r.CurrentRound != roundNum {
		r.mu.Unlock()
		return
	}

	r.Status = StatusRoundReveal
	listing := r.Listings[r.CurrentRound-1]
	r.CurrentListing.ActualPrice = listing.ActualPrice
	r.mu.Unlock()

	r.BroadcastState()

	// Pause for 5 seconds to view reveal, then advance to next round
	time.Sleep(6 * time.Second)
	r.NextRound()
}

// BroadcastState sends full room sync to all clients
func (r *Room) BroadcastState() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.BroadcastFunc == nil {
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"type": "room_state",
		"room": r,
	})
	if err == nil {
		r.BroadcastFunc(data)
	}
}
