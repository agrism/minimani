package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	db "minimani/backend/internal/db/sqlc"
	"minimani/backend/internal/ss"
	"minimani/backend/internal/storage"
)

func main() {
	_ = godotenv.Load("../.env", ".env")

	targetCount := flag.Int("count", 50, "Cik sludinājumus importēt datubāzē")
	flag.Parse()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://minimani:minimani_secret@localhost:5438/minimani_db?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Neizdevās pieslēgties datubāzei: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Datubāze nav sasniedzama: %v", err)
	}

	queries := db.New(pool)

	// Inicializējam S3 glabātuvi
	s3, err := storage.InitStorage()
	if err != nil {
		log.Printf("⚠️ S3 nav aktīvs (%v), tiks izmantotas oriģinālās saites\n", err)
	}

	count, err := ss.ImportHomesCatalog(ctx, pool, queries, s3, *targetCount)
	if err != nil {
		log.Fatalf("Kļūda importā: %v", err)
	}

	totalInDB, _ := queries.CountListings(ctx)
	log.Printf("✨ Imports veiksmīgs! Pievienoti: %d. Kopā datubāzē: %d\n", count, totalInDB)
}
