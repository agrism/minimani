package ss

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DownloadAndSaveImage downloads an image from a URL, saves it locally in mediaDir,
// and returns the local relative URL path (e.g. /media/listings/{uuid}.jpg).
func DownloadAndSaveImage(mediaDir string, rawURL string) (string, error) {
	if rawURL == "" {
		return "", fmt.Errorf("empty image URL")
	}

	// Ensure directory exists
	targetDir := filepath.Join(mediaDir, "listings")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create media directory: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.ss.lv/")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// Determine extension from content-type or URL
	ext := ".jpg"
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "webp") {
		ext = ".webp"
	} else if strings.Contains(contentType, "png") {
		ext = ".png"
	}

	fileName := uuid.New().String() + ext
	filePath := filepath.Join(targetDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write image: %w", err)
	}

	return "/media/listings/" + fileName, nil
}

