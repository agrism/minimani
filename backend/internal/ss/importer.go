package ss

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	internaldb "minimani/backend/internal/db"
	db "minimani/backend/internal/db/sqlc"
	"minimani/backend/internal/storage"
)

// ImportHomesCatalog crawls ss.lv real-estate homes, converts images to WebP in memory,
// uploads them to Hetzner S3, and saves them into the PostgreSQL listings database.
func ImportHomesCatalog(ctx context.Context, pool *pgxpool.Pool, queries *db.Queries, s3 *storage.S3Storage, targetCount int) (int, error) {
	log.Printf("📦 Sākam sludinājumu importu uz datubāzi (mērķis: %d sludinājumi)...\n", targetCount)

	// 1. Nolasām jau esošās saites no datubāzes, lai nedublētu
	existingURLs := make(map[string]bool)
	rows, err := pool.Query(ctx, "SELECT original_url FROM listings WHERE original_url IS NOT NULL")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var u string
			if err := rows.Scan(&u); err == nil && u != "" {
				existingURLs[u] = true
			}
		}
	}
	log.Printf("ℹ️ Datubāzē jau reģistrēti %d unikāli sludinājumi.\n", len(existingURLs))

	maxPages := GetHomesSellMaxPages()
	if maxPages <= 0 {
		maxPages = 91
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffledPages := rng.Perm(maxPages) // Random order of pages 0..maxPages-1

	totalImported := 0

	for _, pIdx := range shuffledPages {
		if ctx.Err() != nil {
			break
		}
		if totalImported >= targetCount {
			break
		}

		pageNum := pIdx + 1
		pageURL := "https://www.ss.lv/lv/real-estate/homes-summer-residences/sell/"
		if pageNum > 1 {
			pageURL = fmt.Sprintf("https://www.ss.lv/lv/real-estate/homes-summer-residences/sell/page%d.html", pageNum)
		}

		links, err := FetchCategoryAdLinks(pageURL, 30)
		if err != nil {
			log.Printf("Kļūda ielādējot lapu %d: %v\n", pageNum, err)
			continue
		}

		for _, link := range links {
			if ctx.Err() != nil {
				break
			}
			if totalImported >= targetCount {
				break
			}
			if existingURLs[link] {
				continue // Jau importēts
			}

			ad, err := ScrapeSSAd(link)
			if err != nil || ad.Price <= 0 || ad.ImageURL == "" {
				continue
			}

			// Pārkodējam uz WebP atmiņā un augšupielādējam S3 (uz diska 0 baiti)
			imagePath := ad.ImageURL
			if s3 != nil {
				s3URL, err := s3.ConvertAndUploadWebP(ctx, ad.ImageURL)
				if err == nil && s3URL != "" {
					imagePath = s3URL
				} else {
					log.Printf("⚠️ S3 augšupielāde neizdevās, izmantojam ss.lv adresi: %v\n", err)
				}
			}

			opts := GeneratePriceOptions(ad.Price)
			optsJSON, _ := json.Marshal(opts)
			packedDesc := PackDescriptionWithParams(ad.Description, ad.Parameters)

			_, err = queries.InsertListing(ctx, db.InsertListingParams{
				Category:    "real_estate",
				Title:       ad.Title,
				Description: pgtype.Text{String: packedDesc, Valid: true},
				ImagePath:   imagePath,
				OriginalUrl: pgtype.Text{String: ad.SourceURL, Valid: true},
				ActualPrice: internaldb.NumericFromFloat(ad.Price),
				Currency:    "EUR",
				Region:      pgtype.Text{String: ad.Region, Valid: ad.Region != ""},
				Options:     optsJSON,
			})
			if err != nil {
				log.Printf("Kļūda saglabājot sludinājumu %s: %v\n", link, err)
				continue
			}

			existingURLs[link] = true
			totalImported++
			log.Printf("  ✅ [%d/%d] %s — %.0f € (S3: %s)\n", totalImported, targetCount, ad.Title, ad.Price, imagePath)

			time.Sleep(200 * time.Millisecond) // Pieklājīga pauze
		}
	}

	log.Printf("🎉 Importēšana pabeigta! Veiksmīgi saglabāti %d jauni sludinājumi.\n", totalImported)
	return totalImported, nil
}

