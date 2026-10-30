package ss

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type ScrapedListing struct {
	Title       string
	Description string
	Price       float64
	Region      string
	ImageURLs   []string
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

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	listing := &ScrapedListing{
		SourceURL: urlStr,
	}

	// 1. Extract Title
	titleRegex := regexp.MustCompile(`<h2[^>]*class="[^"]*head_title[^"]*"[^>]*>(.*?)</h2>`)
	if matches := titleRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
		listing.Title = cleanHTML(matches[1])
	} else {
		// Fallback to title tag
		titleTagRegex := regexp.MustCompile(`<title>(.*?)</title>`)
		if matches := titleTagRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
			listing.Title = strings.TrimSpace(strings.Split(matches[1], "-")[0])
		}
	}

	// 2. Extract Price (e.g. "Cena: 3 850 €" or "3,850 €")
	priceRegex := regexp.MustCompile(`(?:Cena|Price)[:\s]*([0-9\s,\.]+)\s*€`)
	if matches := priceRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
		rawPrice := strings.ReplaceAll(matches[1], " ", "")
		rawPrice = strings.ReplaceAll(rawPrice, ",", ".")
		if p, err := strconv.ParseFloat(rawPrice, 64); err == nil {
			listing.Price = p
		}
	}

	// 3. Extract Description
	descRegex := regexp.MustCompile(`<div[^>]*id="msg_div_msg"[^>]*>(.*?)</div>`)
	if matches := descRegex.FindStringSubmatch(htmlContent); len(matches) > 1 {
		listing.Description = cleanHTML(matches[1])
	}

	// 4. Extract Images (look for gallery thumbs or photo links on i.ss.lv)
	imgRegex := regexp.MustCompile(`https?://[a-zA-Z0-9\.\-_/]+\.ss\.(?:lv|com)/[a-zA-Z0-9\.\-_/]+\.(?:jpg|jpeg|webp)`)
	imgMatches := imgRegex.FindAllString(htmlContent, -1)
	seen := make(map[string]bool)
	for _, m := range imgMatches {
		if !seen[m] && !strings.Contains(m, "logo") && !strings.Contains(m, "icon") {
			seen[m] = true
			listing.ImageURLs = append(listing.ImageURLs, m)
		}
	}

	if listing.Title == "" {
		listing.Title = "ss.lv sludinājums"
	}
	if listing.Region == "" {
		listing.Region = "Latvija"
	}

	_ = doc
	return listing, nil
}

func cleanHTML(s string) string {
	s = regexp.MustCompile(`<br\s*/?>`).ReplaceAllString(s, "\n")
	s = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&amp;", "&")
	return strings.TrimSpace(s)
}
