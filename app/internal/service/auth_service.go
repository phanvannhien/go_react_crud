package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nhienphan/full_rest_app/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic.
type AuthService struct {
	queries   *sqlc.Queries
	jwtSecret string
}

// NewAuthService creates a new AuthService.
func NewAuthService(queries *sqlc.Queries, jwtSecret string) *AuthService {
	return &AuthService{
		queries:   queries,
		jwtSecret: jwtSecret,
	}
}

// AuthUserResponse is the safe user response (no password_hash).
type AuthUserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// RegisterResponse is returned after successful registration.
type RegisterResponse struct {
	User  AuthUserResponse `json:"user"`
	Token string           `json:"token"`
}

// LoginResponse is returned after successful login.
type LoginResponse struct {
	User  AuthUserResponse `json:"user"`
	Token string           `json:"token"`
}

// Register creates a new user with hashed password.
func (s *AuthService) Register(ctx context.Context, email, password, role string) (*RegisterResponse, error) {
	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		IsActive:     true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT
	token, err := s.generateToken(uuidToString(user.ID), user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &RegisterResponse{
		User:  toAuthUserResponse(user.ID, user.Email, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt),
		Token: token,
	}, nil
}

// Login authenticates a user and returns a JWT.
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	// Fetch user by email (includes password_hash for verification)
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT
	token, err := s.generateToken(uuidToString(user.ID), user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		User:  toAuthUserResponse(user.ID, user.Email, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt),
		Token: token,
	}, nil
}

func (s *AuthService) generateToken(userID, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func uuidToString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", id.Bytes[0:4], id.Bytes[4:6], id.Bytes[6:8], id.Bytes[8:10], id.Bytes[10:16])
}

func toAuthUserResponse(id pgtype.UUID, email, role string, isActive bool, createdAt, updatedAt pgtype.Timestamptz) AuthUserResponse {
	return AuthUserResponse{
		ID:        uuidToString(id),
		Email:     email,
		Role:      role,
		IsActive:  isActive,
		CreatedAt: timestampToString(createdAt),
		UpdatedAt: timestampToString(updatedAt),
	}
}

func timestampToString(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return ""
	}
	return ts.Time.Format(time.RFC3339)
}
