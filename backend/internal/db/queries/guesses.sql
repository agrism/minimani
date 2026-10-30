-- name: InsertRoomGuess :one
INSERT INTO room_guesses (
    room_id,
    round_number,
    listing_id,
    player_id,
    guess_value,
    actual_price,
    accuracy_pct,
    score_points,
    time_taken_ms
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
ON CONFLICT (room_id, round_number, player_id) DO UPDATE
SET guess_value = EXCLUDED.guess_value,
    actual_price = EXCLUDED.actual_price,
    accuracy_pct = EXCLUDED.accuracy_pct,
    score_points = EXCLUDED.score_points,
    time_taken_ms = EXCLUDED.time_taken_ms
RETURNING *;

-- name: ListRoundGuesses :many
SELECT g.*, p.nickname, p.avatar
FROM room_guesses g
JOIN room_players p ON g.room_id = p.room_id AND g.player_id = p.player_id
WHERE g.room_id = $1 AND g.round_number = $2
ORDER BY g.score_points DESC, g.time_taken_ms ASC;

-- name: CountRoundGuesses :one
SELECT count(*) FROM room_guesses
WHERE room_id = $1 AND round_number = $2;

