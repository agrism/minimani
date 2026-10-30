-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

-- 1. Sludinājumu banka no ss.lv
CREATE TABLE listings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category        VARCHAR(50) NOT NULL,          -- 'auto', 'real_estate', 'tech', 'pearls'
    title           TEXT NOT NULL,
    description     TEXT,
    image_path      TEXT NOT NULL,                 -- Lokālais /media ceļš uz WebP attēlu
    original_url    TEXT,
    actual_price    NUMERIC(12, 2) NOT NULL,
    currency        VARCHAR(10) NOT NULL DEFAULT 'EUR',
    location        GEOGRAPHY(POINT, 4326),
    region          VARCHAR(100),
    options         JSONB DEFAULT '[]'::jsonb,     -- 4 iepriekš ģenerēti cenu varianti
    times_played    INT NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_listings_active_cat ON listings (category, is_active);
CREATE INDEX idx_listings_location ON listings USING GIST (location);

-- 2. Spēļu istabas (Privātās un Publiskās)
CREATE TABLE rooms (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code              VARCHAR(8) UNIQUE NOT NULL,    -- 4 burtu/ciparu kods, piem. 'AUTO'
    host_player_id    TEXT NOT NULL,
    is_public         BOOLEAN NOT NULL DEFAULT FALSE,
    category          VARCHAR(50) NOT NULL DEFAULT 'all',
    game_mode         VARCHAR(50) NOT NULL DEFAULT 'exact', -- 'exact', 'multiple_choice', 'mix'
    question_count    INT NOT NULL DEFAULT 10,       -- 5, 10, 15, 20
    time_limit_sec    INT NOT NULL DEFAULT 45,
    max_players       INT NOT NULL DEFAULT 8,
    status            VARCHAR(30) NOT NULL DEFAULT 'waiting', -- 'waiting', 'playing', 'finished'
    listing_ids       UUID[] DEFAULT '{}',
    current_round     INT NOT NULL DEFAULT 0,
    round_started_at  TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at       TIMESTAMPTZ
);

CREATE INDEX idx_rooms_public_waiting ON rooms (is_public, status) WHERE is_public AND status = 'waiting';
CREATE INDEX idx_rooms_code ON rooms (code);

-- 3. Istabas spēlētāji
CREATE TABLE room_players (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id         UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    player_id       TEXT NOT NULL,
    nickname        VARCHAR(50) NOT NULL,
    avatar          VARCHAR(20) NOT NULL DEFAULT '🚗',
    is_ready        BOOLEAN NOT NULL DEFAULT FALSE,
    total_score     INT NOT NULL DEFAULT 0,
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_room_player UNIQUE (room_id, player_id)
);

CREATE INDEX idx_room_players_room ON room_players (room_id);

-- 4. Dalībnieku minējumi raundā
CREATE TABLE room_guesses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id         UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    round_number    INT NOT NULL,
    listing_id      UUID NOT NULL REFERENCES listings(id),
    player_id       TEXT NOT NULL,
    guess_value     NUMERIC(12, 2) NOT NULL,
    actual_price    NUMERIC(12, 2) NOT NULL,
    accuracy_pct    NUMERIC(5, 2) NOT NULL,
    score_points    INT NOT NULL DEFAULT 0,
    time_taken_ms   INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_room_round_player UNIQUE (room_id, round_number, player_id)
);

CREATE INDEX idx_room_guesses_room_round ON room_guesses (room_id, round_number);

-- +goose Down
DROP TABLE IF EXISTS room_guesses;
DROP TABLE IF EXISTS room_players;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS listings;

