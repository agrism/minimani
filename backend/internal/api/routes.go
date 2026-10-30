package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"minimani/backend/internal/game"
	"minimani/backend/internal/websocket"
	db "minimani/backend/internal/db/sqlc"
)

type Server struct {
	router  *chi.Mux
	hub     *websocket.Hub
	queries *db.Queries
}

func NewServer(queries *db.Queries, hub *websocket.Hub) *Server {
	s := &Server{
		router:  chi.NewRouter(),
		hub:     hub,
		queries: queries,
	}
	s.setupRoutes()
	return s
}

func (s *Server) Router() http.Handler {
	return s.router
}

func (s *Server) setupRoutes() {
	r := s.router

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS configuration for local development and production
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"app":    "MiniMani Backend",
		})
	})

	// Static media file server
	mediaDir := os.Getenv("MEDIA_DIR")
	if mediaDir == "" {
		mediaDir = "./media"
	}
	_ = os.MkdirAll(mediaDir, 0755)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, mediaDir))
	fileServer(r, "/media", filesDir)

	// API endpoints
	r.Route("/api/rooms", func(r chi.Router) {
		r.Post("/create", s.handleCreateRoom)
		r.Get("/public", s.handleListPublicRooms)
		r.Post("/quick-join", s.handleQuickJoin)
		r.Get("/{code}", s.handleGetRoom)
	})

	// WebSocket route
	r.Get("/ws/rooms/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		s.hub.HandleWebSocket(w, r, code)
	})
}

// Generate fun Latvian 4-letter room codes
var roomCodeWords = []string{
	"AUTO", "CENA", "RIGA", "MAJA", "LATS", "ZELTS", "GOLF", "VECS", "JAUNS", "PIRKT",
	"PARD", "ALUS", "KRAV", "ZIVE", "RUPI", "E46", "AUDI", "VOLV", "OPEL", "SS",
}

func generateRoomCode() string {
	word := roomCodeWords[rand.Intn(len(roomCodeWords))]
	if len(word) < 4 {
		return fmt.Sprintf("%s%d", word, rand.Intn(90)+10)
	}
	return word
}

type CreateRoomRequest struct {
	HostID   string `json:"host_id"`
	IsPublic bool   `json:"is_public"`
	Category string `json:"category"`
	GameMode string `json:"game_mode"`
}

func (s *Server) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	code := generateRoomCode()
	room := s.hub.GetOrCreateRoom(code, req.HostID, req.IsPublic)
	if req.Category != "" {
		room.Category = req.Category
	}
	if req.GameMode != "" {
		room.GameMode = req.GameMode
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code": code,
		"room": room,
	})
}

func (s *Server) handleListPublicRooms(w http.ResponseWriter, r *http.Request) {
	rooms := s.hub.ListPublicRooms()

	type PublicRoomItem struct {
		Code          string `json:"code"`
		Category      string `json:"category"`
		GameMode      string `json:"game_mode"`
		PlayerCount   int    `json:"player_count"`
		MaxPlayers    int    `json:"max_players"`
		QuestionCount int    `json:"question_count"`
		Status        string `json:"status"`
	}

	var list []PublicRoomItem
	for _, room := range rooms {
		list = append(list, PublicRoomItem{
			Code:          room.Code,
			Category:      room.Category,
			GameMode:      room.GameMode,
			PlayerCount:   len(room.Players),
			MaxPlayers:    room.MaxPlayers,
			QuestionCount: room.QuestionCount,
			Status:        string(room.Status),
		})
	}

	if list == nil {
		list = []PublicRoomItem{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Server) handleQuickJoin(w http.ResponseWriter, r *http.Request) {
	rooms := s.hub.ListPublicRooms()

	// Find first room with room for more players
	for _, room := range rooms {
		if len(room.Players) < room.MaxPlayers && room.Status == game.StatusLobby {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"code": room.Code,
			})
			return
		}
	}

	// If no public rooms available, create one
	code := generateRoomCode()
	_ = s.hub.GetOrCreateRoom(code, "system", true)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code": code,
	})
}

func (s *Server) handleGetRoom(w http.ResponseWriter, r *http.Request) {
	code := strings.ToUpper(chi.URLParam(r, "code"))
	room := s.hub.GetRoom(code)
	if room == nil {
		http.Error(w, `{"error":"Istaba nav atrasta"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(room)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	handler := func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	}

	r.Get(path, handler)
	r.Head(path, handler)
}

