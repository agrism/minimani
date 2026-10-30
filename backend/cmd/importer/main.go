package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	internaldb "minimani/backend/internal/db"
	db "minimani/backend/internal/db/sqlc"
	"minimani/backend/internal/ss"
)

func main() {
	_ = godotenv.Load("../.env", ".env")

	action := flag.String("action", "scrape_live", "Action to run: 'scrape_live' or 'seed'")
	limitPerCat := flag.Int("limit", 6, "Number of ads to scrape per category")
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

	queries := db.New(pool)

	if *action == "scrape_live" {
		scrapeLiveFromSS(ctx, queries, *limitPerCat)
	} else if *action == "seed" {
		seedDatabase(ctx, queries)
	} else {
		fmt.Printf("Nezināma darbība: %s\n", *action)
	}
}

func scrapeLiveFromSS(ctx context.Context, queries *db.Queries, limitPerCat int) {
	log.Println("🕷️  Sākam 100% REĀLU sludinājumu un bilžu ievākšanu no ss.lv...")

	mediaDir := os.Getenv("MEDIA_DIR")
	if mediaDir == "" {
		mediaDir = "./media"
	}

	categories := []struct {
		ID  string
		URL string
	}{
		{"auto", "https://www.ss.lv/lv/transport/cars/today/"},
		{"real_estate", "https://www.ss.lv/lv/real-estate/flats/riga/today/"},
		{"real_estate", "https://www.ss.lv/lv/real-estate/homes-summer-residences/today/"},
		{"tech", "https://www.ss.lv/lv/agriculture/agricultural-machinery/tractors/"},
		{"tech", "https://www.ss.lv/lv/construction/garden-technics/mowers-riders/"},
		{"tech", "https://www.ss.lv/lv/electronics/telephones/mobile-phones/apple/"},
		{"pearls", "https://www.ss.lv/lv/home-stuff/antique/furniture/"},
		{"pearls", "https://www.ss.lv/lv/home-stuff/furniture-interior/carpets/"},
	}

	totalSaved := 0

	for _, cat := range categories {
		log.Printf("📂 Nolasām kategoriju: %s (%s)...\n", cat.ID, cat.URL)
		links, err := ss.FetchCategoryAdLinks(cat.URL, limitPerCat)
		if err != nil {
			log.Printf("Kļūda nolasot saites no %s: %v", cat.URL, err)
			continue
		}

		for _, link := range links {
			time.Sleep(300 * time.Millisecond) // Respectful crawl delay

			ad, err := ss.ScrapeSSAd(link)
			if err != nil || ad.Price <= 0 || ad.ImageURL == "" {
				continue
			}

			// Lejupielādējam reālo attēlu no i.ss.lv uz lokālo /media mapi
			localImagePath, err := ss.DownloadAndSaveImage(mediaDir, ad.ImageURL)
			if err != nil {
				log.Printf("Neizdevās lejupielādēt attēlu priekš %s: %v", ad.Title, err)
				localImagePath = ad.ImageURL // fallback to direct URL
			}

			optionsJSON := ss.GenerateOptionsJSON(ad.Price)

			_, err = queries.InsertListing(ctx, db.InsertListingParams{
				Category:    cat.ID,
				Title:       ad.Title,
				Description: pgtype.Text{String: ad.Description, Valid: ad.Description != ""},
				ImagePath:   localImagePath,
				OriginalUrl: pgtype.Text{String: ad.SourceURL, Valid: true},
				ActualPrice: internaldb.NumericFromFloat(ad.Price),
				Currency:    "EUR",
				Region:      pgtype.Text{String: ad.Region, Valid: ad.Region != ""},
				Options:     optionsJSON,
			})
			if err != nil {
				log.Printf("Kļūda saglabājot DB: %v", err)
				continue
			}

			log.Printf("  ✓ [%s] %s — %.2f € (bilde: %s)\n", cat.ID, ad.Title, ad.Price, localImagePath)
			totalSaved++
		}
	}

	count, _ := queries.CountListings(ctx)
	log.Printf("🎉 Gatavs! Pievienoti %d 100%% reāli ss.lv sludinājumi ar īstām bildēm! Kopā DB: %d\n", totalSaved, count)
}

func seedDatabase(ctx context.Context, queries *db.Queries) {
	log.Println("🌱 Sēklas imports...")
}
