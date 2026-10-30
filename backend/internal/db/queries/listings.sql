-- name: GetListingByID :one
SELECT * FROM listings
WHERE id = $1 AND is_active = TRUE;

-- name: GetRandomListings :many
SELECT * FROM listings
WHERE is_active = TRUE
  AND (category = $1 OR $1 = 'all')
ORDER BY RANDOM()
LIMIT $2;

-- name: InsertListing :one
INSERT INTO listings (
    category,
    title,
    description,
    image_path,
    original_url,
    actual_price,
    currency,
    region,
    options
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: IncrementListingTimesPlayed :exec
UPDATE listings
SET times_played = times_played + 1
WHERE id = $1;

-- name: CountListings :one
SELECT count(*) FROM listings WHERE is_active = TRUE;

