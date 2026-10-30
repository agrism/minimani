package auth

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	GamesPlayed int       `json:"games_played"`
	GamesWon    int       `json:"games_won"`
	TotalScore  int       `json:"total_score"`
	CreatedAt   time.Time `json:"created_at"`
}

type AuthService struct {
	pool      *pgxpool.Pool
	jwtSecret []byte
}

func NewAuthService(pool *pgxpool.Pool) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "minimani-secret-jwt-key-2026-super-safe"
	}
	return &AuthService{
		pool:      pool,
		jwtSecret: []byte(secret),
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, nickname, avatar string) (*User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	nickname = strings.TrimSpace(nickname)

	// Validation
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, "", errors.New("nederīgs e-pasta formāts")
	}
	if len(password) < 6 {
		return nil, "", errors.New("parolei jābūt vismaz 6 simbolus garai")
	}
	if len(nickname) < 2 || len(nickname) > 30 {
		return nil, "", errors.New("segvārdam (nikam) jābūt no 2 līdz 30 simboliem")
	}
	if avatar == "" {
		avatar = "🚗"
	}

	// Check if email already exists
	var exists bool
	err := s.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE lower(email) = $1)", email).Scan(&exists)
	if err != nil {
		return nil, "", fmt.Errorf("datubāzes kļūda: %w", err)
	}
	if exists {
		return nil, "", errors.New("lietotājs ar šādu e-pastu jau eksistē")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("paroles apstrādes kļūda: %w", err)
	}

	// Insert user
	var user User
	err = s.pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, nickname, avatar)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, nickname, avatar, games_played, games_won, total_score, created_at
	`, email, string(hash), nickname, avatar).Scan(
		&user.ID, &user.Email, &user.Nickname, &user.Avatar,
		&user.GamesPlayed, &user.GamesWon, &user.TotalScore, &user.CreatedAt,
	)
	if err != nil {
		return nil, "", fmt.Errorf("neizdevās reģistrēt lietotāju: %w", err)
	}

	// Generate token
	token, err := s.generateToken(user.ID, user.Email, user.Nickname)
	if err != nil {
		return nil, "", fmt.Errorf("neizdevās izveidot sesijas marķieri: %w", err)
	}

	return &user, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	var user User
	var hash string
	err := s.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, nickname, avatar, games_played, games_won, total_score, created_at
		FROM users
		WHERE lower(email) = $1
	`, email).Scan(
		&user.ID, &user.Email, &hash, &user.Nickname, &user.Avatar,
		&user.GamesPlayed, &user.GamesWon, &user.TotalScore, &user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", errors.New("nepareizs e-pasts vai parole")
		}
		return nil, "", fmt.Errorf("datubāzes kļūda: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, "", errors.New("nepareizs e-pasts vai parole")
	}

	token, err := s.generateToken(user.ID, user.Email, user.Nickname)
	if err != nil {
		return nil, "", fmt.Errorf("neizdevās izveidot sesijas marķieri: %w", err)
	}

	return &user, token, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := s.pool.QueryRow(ctx, `
		SELECT id, email, nickname, avatar, games_played, games_won, total_score, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.Nickname, &user.Avatar,
		&user.GamesPlayed, &user.GamesWon, &user.TotalScore, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) VerifyToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("nederīgs vai beidzies sesijas marķieris")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("nederīgas marķiera prasības")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", errors.New("lietotāja ID nav atrasts marķierī")
	}
	return userID, nil
}

func (s *AuthService) generateToken(userID, email, nickname string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"email":    email,
		"nickname": nickname,
		"exp":      time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days
		"iat":      time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.jwtSecret)
}
