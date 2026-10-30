package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	webp "github.com/mayahiro/go-webp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	publicBase string
}

var GlobalStorage *S3Storage

// InitStorage initializes S3 client from environment variables
func InitStorage() (*S3Storage, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		endpoint = "hel1.your-objectstorage.com"
	}
	// Strip https:// or http:// prefix for minio-go client
	rawHost := strings.TrimPrefix(endpoint, "https://")
	rawHost = strings.TrimPrefix(rawHost, "http://")
	rawHost = strings.TrimRight(rawHost, "/")

	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucketName := os.Getenv("S3_BUCKET")
	if bucketName == "" {
		bucketName = "minimani"
	}

	if accessKey == "" || secretKey == "" {
		log.Println("⚠️ S3_ACCESS_KEY vai S3_SECRET_KEY nav norādīts, S3 glabātuve nav aktīva")
		return nil, fmt.Errorf("missing S3 credentials")
	}

	client, err := minio.New(rawHost, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init S3 client: %w", err)
	}

	// Format public base URL: https://minimani.hel1.your-objectstorage.com or https://hel1.your-objectstorage.com/minimani
	publicBase := fmt.Sprintf("https://%s.%s", bucketName, rawHost)

	s3 := &S3Storage{
		client:     client,
		bucketName: bucketName,
		endpoint:   rawHost,
		publicBase: publicBase,
	}

	GlobalStorage = s3
	log.Printf("✅ S3 glabātuve inicializēta: bucket=%s endpoint=%s publicBase=%s\n", bucketName, rawHost, publicBase)
	return s3, nil
}

// ConvertAndUploadWebP downloads an image from rawURL into memory, converts it to WebP in memory,
// and uploads it directly to Hetzner S3 bucket. Nothing is written to disk.
func (s *S3Storage) ConvertAndUploadWebP(ctx context.Context, rawURL string) (string, error) {
	if rawURL == "" {
		return "", fmt.Errorf("empty raw image URL")
	}

	// 1. Download image to memory buffer
	httpClient := &http.Client{Timeout: 12 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("create req error: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.ss.lv/")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("image fetch error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("image fetch status %d", resp.StatusCode)
	}

	rawBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read image bytes error: %w", err)
	}

	// 2. Decode image in memory
	img, _, err := image.Decode(bytes.NewReader(rawBytes))
	if err != nil {
		return "", fmt.Errorf("image decode error: %w", err)
	}

	// 3. Encode image to WebP format in memory
	var webpBuf bytes.Buffer
	err = webp.Encode(&webpBuf, img, &webp.Options{
		Compression: webp.CompressionLossy,
		Quality:     82,
	})
	if err != nil {
		return "", fmt.Errorf("webp encode error: %w", err)
	}

	// 4. Upload to S3 bucket in memory
	objectKey := fmt.Sprintf("listings/%s.webp", uuid.New().String())
	contentLength := int64(webpBuf.Len())

	_, err = s.client.PutObject(ctx, s.bucketName, objectKey, &webpBuf, contentLength, minio.PutObjectOptions{
		ContentType: "image/webp",
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	})
	if err != nil {
		return "", fmt.Errorf("s3 put object error: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s", s.publicBase, objectKey)
	return publicURL, nil
}

