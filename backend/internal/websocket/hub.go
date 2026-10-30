package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"minimani/backend/internal/game"
	db "minimani/backend/internal/db/sqlc"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	roomCode string
	playerID string
	send     chan []byte
}

type Hub struct {
	mu          sync.RWMutex
	rooms       map[string]*game.Room
	roomClients map[string]map[*Client]bool
	queries     *db.Queries
}

func NewHub(queries *db.Queries) *Hub {
	return &Hub{
		rooms:       make(map[string]*game.Room),
		roomClients: make(map[string]map[*Client]bool),
		queries:     queries,
	}
}

// GetOrCreateRoom gets an existing room or creates a new one
func (h *Hub) GetOrCreateRoom(code string, hostID string, isPublic bool) *game.Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	code = strings.ToUpper(code)
	if room, exists := h.rooms[code]; exists {
		return room
	}

	room := game.NewRoom(code, hostID, isPublic, h.queries, func(msg []byte) {
		h.BroadcastToRoom(code, msg)
	})
	h.rooms[code] = room
	h.roomClients[code] = make(map[*Client]bool)
	return room
}

func (h *Hub) GetRoom(code string) *game.Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[strings.ToUpper(code)]
}

// ListPublicRooms returns list of active public rooms in waiting status
func (h *Hub) ListPublicRooms() []*game.Room {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var list []*game.Room
	for _, r := range h.rooms {
		if r.IsPublic && (r.Status == game.StatusLobby) {
			list = append(list, r)
		}
	}
	return list
}

func (h *Hub) BroadcastToRoom(roomCode string, msg []byte) {
	h.mu.RLock()
	clients := h.roomClients[strings.ToUpper(roomCode)]
	h.mu.RUnlock()

	for client := range clients {
		select {
		case client.send <- msg:
		default:
			// client channel blocked
		}
	}
}

// HandleWebSocket handles incoming client connection
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, roomCode string) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Printf("WebSocket accept error: %v", err)
		return
	}

	roomCode = strings.ToUpper(roomCode)
	room := h.GetRoom(roomCode)
	if room == nil {
		_ = conn.Close(websocket.StatusPolicyViolation, "Istaba neeksistē")
		return
	}

	client := &Client{
		hub:      h,
		conn:     conn,
		roomCode: roomCode,
		send:     make(chan []byte, 64),
	}

	h.mu.Lock()
	if _, exists := h.roomClients[roomCode]; !exists {
		h.roomClients[roomCode] = make(map[*Client]bool)
	}
	h.roomClients[roomCode][client] = true
	h.mu.Unlock()

	// Write pump
	go func() {
		defer func() {
			_ = conn.Close(websocket.StatusNormalClosure, "")
		}()
		for msg := range client.send {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := conn.Write(ctx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				return
			}
		}
	}()

	// Send current room state immediately on connect
	room.BroadcastState()

	// Read pump
	defer func() {
		h.mu.Lock()
		delete(h.roomClients[roomCode], client)
		if client.playerID != "" {
			room.RemovePlayer(client.playerID)
		}
		h.mu.Unlock()
		close(client.send)
		room.BroadcastState()
	}()

	for {
		_, msg, err := conn.Read(r.Context())
		if err != nil {
			break
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(msg, &payload); err != nil {
			continue
		}

		action, _ := payload["action"].(string)
		switch action {
		case "join":
			playerID, _ := payload["player_id"].(string)
			nickname, _ := payload["nickname"].(string)
			avatar, _ := payload["avatar"].(string)
			client.playerID = playerID
			room.AddPlayer(playerID, nickname, avatar)
			room.BroadcastState()

		case "vote":
			playerID, _ := payload["player_id"].(string)
			var count int
			if v, ok := payload["question_count"].(float64); ok {
				count = int(v)
			}
			mode, _ := payload["game_mode"].(string)
			cat, _ := payload["category"].(string)
			room.Vote(playerID, count, mode, cat)
			room.BroadcastState()

		case "ready":
			playerID, _ := payload["player_id"].(string)
			ready, _ := payload["is_ready"].(bool)
			allReady := room.SetReady(playerID, ready)
			room.BroadcastState()

			// Auto start if all ready
			if allReady && room.Status == game.StatusLobby {
				_ = room.StartGame(context.Background())
			}

		case "start":
			if room.Status == game.StatusLobby {
				_ = room.StartGame(context.Background())
			}

		case "guess":
			playerID, _ := payload["player_id"].(string)
			var price float64
			if v, ok := payload["guess_price"].(float64); ok {
				price = v
			}
			room.SubmitGuess(playerID, price)
		}
	}
}

