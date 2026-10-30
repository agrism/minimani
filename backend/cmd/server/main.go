package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"minimani/backend/internal/api"
	db "minimani/backend/internal/db/sqlc"
	"minimani/backend/internal/websocket"
)

func main() {
	_ = godotenv.Load("../.env", ".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://minimani:minimani_secret@localhost:5438/minimani_db?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Neizdevās izveidot DB pieslēgumu: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Datubāze nav sasniedzama: %v", err)
	}
	log.Println("✅ PostgreSQL (PostGIS) veiksmīgi pieslēgts!")

	queries := db.New(pool)
	hub := websocket.NewHub(queries)
	server := api.NewServer(queries, hub)

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      server.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 MiniMani backend serveris startēts uz http://localhost:%s\n", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Servera kļūda: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Serveris tiek apturēts...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Kļūda apturot serveri: %v", err)
	}
	log.Println("Serveris droši apturēts.")
}

