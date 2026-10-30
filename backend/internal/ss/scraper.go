package ss

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"minimani/backend/internal/storage"
)

type ParameterItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type StoredListingDetails struct {
	Desc   string          `json:"desc"`
	Params []ParameterItem `json:"params"`
}

type ScrapedListing struct {
	Title       string
	Description string
	Price       float64
	Region      string
	ImageURL    string
	SourceURL   string
	Parameters  []ParameterItem
}

// PackDescriptionWithParams serializes description and parameters to JSON for storage
func PackDescriptionWithParams(desc string, params []ParameterItem) string {
	details := StoredListingDetails{
		Desc:   desc,
		Params: params,
	}
	b, err := json.Marshal(details)
	if err != nil {
		return desc
	}
	return string(b)
}

// UnpackDescriptionWithParams decodes description and parameters from stored string
func UnpackDescriptionWithParams(raw string) (string, []ParameterItem) {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "{") {
		var details StoredListingDetails
		if err := json.Unmarshal([]byte(raw), &details); err == nil && len(details.Params) > 0 {
			return details.Desc, details.Params
		}
	}
	return raw, nil
}

// ExtractParameters extracts the key-value attributes from the options_list table
func ExtractParameters(htmlContent string) []ParameterItem {
	var params []ParameterItem

	optRegex := regexp.MustCompile(`(?s)<td[^>]*class=["']ads_opt_name[^"']*["'][^>]*>(.*?)</td>\s*<td[^>]*class=["']ads_opt[^"']*["'][^>]*>(.*?)</td>`)
	matches := optRegex.FindAllStringSubmatch(htmlContent, -1)

	for _, m := range matches {
		label := cleanHTML(m[1])
		label = strings.TrimSuffix(label, ":")
		label = strings.TrimSpace(label)

		val := cleanHTML(m[2])
		val = strings.ReplaceAll(val, "[Karte]", "")
		val = strings.TrimSpace(val)

		// Never leak price!
		lowerLabel := strings.ToLower(label)
		if strings.Contains(lowerLabel, "cena") || strings.Contains(lowerLabel, "price") || strings.Contains(lowerLabel, "vērtība") {
			continue
		}

		if label != "" && val != "" {
			params = append(params, ParameterItem{
				Label: label,
				Value: val,
			})
		}
	}
	return params
}

// ScrapeSSAd parses an ss.lv advertisement URL and extracts metadata
func ScrapeSSAd(urlStr string) (*ScrapedListing, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "lv,en;q=0.9")
	req.Header.Set("Referer", "https://www.ss.lv/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pieprasījuma kļūda: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ss.lv atgrieza statusu: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	htmlContent := string(bodyBytes)

	listing := &ScrapedListing{
		SourceURL: urlStr,
	}

	// 1. Extract Price from JS variable or HTML table
	// e.g. var MSG_PRICE = 16490.00;
	msgPriceRegex := regexp.MustCompile(`var\s+MSG_PRICE\s*=\s*([0-9\.]+);`)
	if matches := msgPriceRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
		if p, err := strconv.ParseFloat(matches[1], 64); err == nil {
			listing.Price = p
		}
	}

	if listing.Price == 0 {
		priceRegex := regexp.MustCompile(`(?:Cena|Price)[:\s]*([0-9\s,\.]+)\s*€`)
		if matches := priceRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
			rawPrice := strings.ReplaceAll(matches[1], " ", "")
			rawPrice = strings.ReplaceAll(rawPrice, ",", ".")
			if p, err := strconv.ParseFloat(rawPrice, 64); err == nil {
				listing.Price = p
			}
		}
	}

	// 2. Extract Title
	titleRegex := regexp.MustCompile(`<h2[^>]*class="[^"]*head_title[^"]*"[^>]*>(.*?)</h2>`)
	if matches := titleRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
		listing.Title = cleanHTML(matches[1])
	} else {
		titleTagRegex := regexp.MustCompile(`<title>(.*?)</title>`)
		if matches := titleTagRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
			raw := strings.Split(matches[1], "-")[0]
			listing.Title = strings.TrimSpace(raw)
		}
	}

	// 3. Extract Description
	descRegex := regexp.MustCompile(`(?s)<div[^>]*id="msg_div_msg"[^>]*>(.*?)</div>`)
	if matches := descRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
		listing.Description = cleanHTML(matches[1])
	}

	// 4. Extract Real 800px Image from i.ss.lv
	// e.g. https://i.ss.lv/gallery/8/1563/390717/volvo-v60-78143308.800.jpg
	img800Regex := regexp.MustCompile(`https?://i\.ss\.(?:lv|com)/gallery/[a-zA-Z0-9\.\-_/]+\.800\.jpg`)
	if m := img800Regex.FindString(htmlContent); m != "" {
		listing.ImageURL = m
	} else {
		// Fallback to any gallery image
		imgAnyRegex := regexp.MustCompile(`https?://i\.ss\.(?:lv|com)/gallery/[a-zA-Z0-9\.\-_/]+\.(?:jpg|jpeg|webp)`)
		if m := imgAnyRegex.FindString(htmlContent); m != "" {
			listing.ImageURL = m
		}
	}

	// 5. Extract Table Parameters (Platība, Istabas, Stāvu skaits, Zemes platība, etc.)
	listing.Parameters = ExtractParameters(htmlContent)

	// 6. Extract Region / City
	regionRegex := regexp.MustCompile(`(?:Vieta|Pilsēta|Rajons)[:\s]*<[^>]*>([^<]+)<`)
	if matches := regionRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
		listing.Region = cleanHTML(matches[1])
	}
	if listing.Region == "" || listing.Region == "Latvija" {
		for _, p := range listing.Parameters {
			if p.Label == "Pilsēta" || p.Label == "Pilsēta/pagasts" || p.Label == "Rajons" {
				listing.Region = p.Value
				break
			}
		}
	}
	if listing.Region == "" {
		listing.Region = "Latvija"
	}

	if listing.Title == "" {
		listing.Title = "Mājas, vasarnīcas"
	}
	// Sanitize title and description to remove price spoilers
	listing.Title = SanitizeListingText(listing.Title)
	listing.Description = SanitizeListingText(listing.Description)

	return listing, nil
}

// SanitizeListingText strips price mentions and ss.lv prefixes so players cannot cheat from the title or description
func SanitizeListingText(s string) string {
	s = cleanHTML(s)
	// Remove leading SS.LV prefix
	s = regexp.MustCompile(`(?i)^\s*ss\.(?:lv|com)\s*[-:]*\s*`).ReplaceAllString(s, "")
	// Remove "Cena 12 500 €" or "Cena: 500 €"
	s = regexp.MustCompile(`(?i)(?:cena|price|vērtība)[:\s]*[0-9\s,\.]+\s*(?:€|eur|eiro)?`).ReplaceAllString(s, "")
	// Remove standalone price patterns like "10 850 €"
	s = regexp.MustCompile(`[0-9][0-9\s,\.]*\s*(?:€|EUR|eur)`).ReplaceAllString(s, "")
	// Clean extra commas or dots left over
	s = regexp.MustCompile(`\s*[,.]\s*[,.]+\s*`).ReplaceAllString(s, ". ")
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// FetchCategoryAdLinks crawls an ss.lv category page and extracts ad links
func FetchCategoryAdLinks(categoryURL string, limit int) ([]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", categoryURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.ss.lv/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	linkRegex := regexp.MustCompile(`href="(/msg/lv/[^"]+\.html)"`)
	matches := linkRegex.FindAllStringSubmatch(string(bodyBytes), -1)

	var links []string
	seen := make(map[string]bool)
	for _, m := range matches {
		fullURL := "https://www.ss.lv" + m[1]
		if !seen[fullURL] {
			seen[fullURL] = true
			links = append(links, fullURL)
			if len(links) >= limit {
				break
			}
		}
	}

	return links, nil
}

func cleanHTML(s string) string {
	s = regexp.MustCompile(`<br\s*/?>`).ReplaceAllString(s, " ")
	s = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "\t", " ")
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

var (
	cachedMaxPages     int
	cachedMaxPagesTime time.Time
)

// GetHomesSellMaxPages dynamically determines total page count from "Iepriekšējie" link (e.g. 91)
func GetHomesSellMaxPages() int {
	if cachedMaxPages > 0 && time.Since(cachedMaxPagesTime) < 30*time.Minute {
		return cachedMaxPages
	}

	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest("GET", "https://www.ss.lv/lv/real-estate/homes-summer-residences/sell/", nil)
	if err != nil {
		return 91
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.ss.lv/")

	resp, err := client.Do(req)
	if err != nil {
		if cachedMaxPages > 0 {
			return cachedMaxPages
		}
		return 91
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	htmlStr := string(bodyBytes)

	// Look for rel="prev" or "[Iepriekšējie]" link: e.g. href="/lv/real-estate/homes-summer-residences/sell/page91.html"
	re := regexp.MustCompile(`href=["'](?:/lv/real-estate/homes-summer-residences/sell/)?page([0-9]+)\.html["'][^>]*>(?:<img[^>]*>)?\s*Iepriekšējie`)
	if m := re.FindStringSubmatch(htmlStr); len(m) > 1 {
		if p, err := strconv.Atoi(m[1]); err == nil && p > 0 {
			cachedMaxPages = p
			cachedMaxPagesTime = time.Now()
			return p
		}
	}

	// Fallback regex looking for rel="prev"
	reFallback := regexp.MustCompile(`rel=["']prev["'][^>]*href=["'][^"']*page([0-9]+)\.html["']`)
	if m := reFallback.FindStringSubmatch(htmlStr); len(m) > 1 {
		if p, err := strconv.Atoi(m[1]); err == nil && p > 0 {
			cachedMaxPages = p
			cachedMaxPagesTime = time.Now()
			return p
		}
	}

	if cachedMaxPages > 0 {
		return cachedMaxPages
	}
	return 91
}

// FetchLiveListings crawls ss.lv live for real-estate homes by picking random pages from 1 to maxPages
func FetchLiveListings(ctx context.Context, mediaDir string, count int, excludeURLs map[string]bool) ([]ScrapedListing, error) {
	maxPages := GetHomesSellMaxPages()
	if maxPages <= 0 {
		maxPages = 91
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Pick random unique pages out of 1..maxPages
	numPagesToFetch := (count / 4) + 3 // e.g. 4-6 random pages from 91
	if numPagesToFetch > maxPages {
		numPagesToFetch = maxPages
	}

	pageSet := make(map[int]bool)
	var randomPages []int
	for len(randomPages) < numPagesToFetch {
		p := rng.Intn(maxPages) + 1
		if !pageSet[p] {
			pageSet[p] = true
			randomPages = append(randomPages, p)
		}
	}

	var allLinks []string
	seen := make(map[string]bool)

	for _, pageNum := range randomPages {
		if ctx.Err() != nil {
			break
		}
		pageURL := "https://www.ss.lv/lv/real-estate/homes-summer-residences/sell/"
		if pageNum > 1 {
			pageURL = fmt.Sprintf("https://www.ss.lv/lv/real-estate/homes-summer-residences/sell/page%d.html", pageNum)
		}

		links, err := FetchCategoryAdLinks(pageURL, 30)
		if err == nil {
			for _, l := range links {
				if !seen[l] && !excludeURLs[l] {
					seen[l] = true
					allLinks = append(allLinks, l)
				}
			}
		}
		if len(allLinks) >= count*4 {
			break
		}
	}

	if len(allLinks) == 0 {
		return nil, fmt.Errorf("neizdevās atrast jaunas sludinājumu saites ss.lv")
	}

	// Nejauši sajaucam saites (random shuffle)
	rng.Shuffle(len(allLinks), func(i, j int) {
		allLinks[i], allLinks[j] = allLinks[j], allLinks[i]
	})

	var results []ScrapedListing
	for _, link := range allLinks {
		if ctx.Err() != nil {
			break
		}
		ad, err := ScrapeSSAd(link)
		if err != nil || ad.Price <= 0 || ad.ImageURL == "" {
			continue
		}

		// Pārveidojam uz WebP atmiņā un augšupielādējam tieši uz Hetzner S3 (uz diska 0 baiti!)
		if storage.GlobalStorage != nil {
			s3URL, err := storage.GlobalStorage.ConvertAndUploadWebP(ctx, ad.ImageURL)
			if err == nil && s3URL != "" {
				ad.ImageURL = s3URL
			} else {
				log.Printf("⚠️ S3 WebP augšupielāde neizdevās, izmantojam ss.lv adresi: %v", err)
			}
		}

		results = append(results, *ad)
		if len(results) >= count {
			break
		}
		time.Sleep(120 * time.Millisecond) // Respectful crawl delay
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("neizdevās nolasīt derīgus sludinājumus")
	}

	return results, nil
}
