package ss

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ScrapedListing struct {
	Title       string
	Description string
	Price       float64
	Region      string
	ImageURL    string
	SourceURL   string
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

	// 5. Extract Region / City
	regionRegex := regexp.MustCompile(`(?:Vieta|Pilsēta|Rajons)[:\s]*<[^>]*>([^<]+)<`)
	if matches := regionRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
		listing.Region = cleanHTML(matches[1])
	}
	if listing.Region == "" {
		listing.Region = "Latvija"
	}

	if listing.Title == "" {
		listing.Title = "ss.lv sludinājums"
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
