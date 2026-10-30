-- name: CreateRoom :one
INSERT INTO rooms (
    code,
    host_player_id,
    is_public,
    category,
    game_mode,
    question_count,
    time_limit_sec,
    max_players,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, 'waiting'
)
RETURNING *;

-- name: GetRoomByCode :one
SELECT * FROM rooms
WHERE code = $1;

-- name: GetRoomByID :one
SELECT * FROM rooms
WHERE id = $1;

-- name: ListPublicWaitingRooms :many
SELECT r.*, count(p.id)::int as player_count
FROM rooms r
LEFT JOIN room_players p ON r.id = p.room_id
WHERE r.is_public = TRUE AND r.status = 'waiting'
GROUP BY r.id
ORDER BY r.created_at DESC
LIMIT 20;

-- name: UpdateRoomSettings :one
UPDATE rooms
SET question_count = $2,
    game_mode = $3,
    category = $4
WHERE id = $1
RETURNING *;

-- name: StartRoomGame :one
UPDATE rooms
SET status = 'playing',
    listing_ids = $2,
    current_round = 1,
    round_started_at = now()
WHERE id = $1
RETURNING *;

-- name: AdvanceRoomRound :one
UPDATE rooms
SET current_round = $2,
    round_started_at = now()
WHERE id = $1
RETURNING *;

-- name: FinishRoomGame :one
UPDATE rooms
SET status = 'finished',
    finished_at = now()
WHERE id = $1
RETURNING *;

-- name: UpsertRoomPlayer :one
INSERT INTO room_players (
    room_id,
    player_id,
    nickname,
    avatar,
    is_ready,
    total_score
) VALUES (
    $1, $2, $3, $4, $5, $6
)
ON CONFLICT (room_id, player_id) DO UPDATE
SET nickname = EXCLUDED.nickname,
    avatar = EXCLUDED.avatar,
    is_ready = EXCLUDED.is_ready,
    total_score = EXCLUDED.total_score
RETURNING *;

-- name: ListRoomPlayers :many
SELECT * FROM room_players
WHERE room_id = $1
ORDER BY total_score DESC, joined_at ASC;

-- name: SetPlayerReady :exec
UPDATE room_players
SET is_ready = $3
WHERE room_id = $1 AND player_id = $2;

-- name: AddPlayerScore :exec
UPDATE room_players
SET total_score = total_score + $3
WHERE room_id = $1 AND player_id = $2;

-- name: RemoveRoomPlayer :exec
DELETE FROM room_players
WHERE room_id = $1 AND player_id = $2;

