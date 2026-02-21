package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nhienphan/full_rest_app/db/sqlc"
)

// UserService handles user business logic.
type UserService struct {
	queries *sqlc.Queries
}

// NewUserService creates a new UserService.
func NewUserService(queries *sqlc.Queries) *UserService {
	return &UserService{queries: queries}
}

// UserResponse is the safe user response (no password_hash).
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UserListResult holds paginated user results.
type UserListResult struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}

// UserListParams holds list query params.
type UserListParams struct {
	Search    *string
	Role      *string
	IsActive  *bool
	SortField string
	SortOrder string
	Page      int
	Limit     int
}

// AllowedSortFields for users.
var userAllowedSortFields = map[string]bool{
	"email":      true,
	"role":       true,
	"created_at": true,
}

// ListUsers returns a paginated list of users.
func (s *UserService) ListUsers(ctx context.Context, params UserListParams) (*UserListResult, error) {
	// Validate sort field
	if !userAllowedSortFields[params.SortField] {
		params.SortField = "created_at"
	}
	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		params.SortOrder = "desc"
	}

	// Enforce limit cap
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	offset := (params.Page - 1) * params.Limit

	// Build filter params - use empty/zero for nil (sqlc NULL pattern)
	var searchParam string
	var roleParam string
	var isActiveParam bool

	searchNull := true
	roleNull := true
	isActiveNull := true

	if params.Search != nil && *params.Search != "" {
		searchParam = *params.Search
		searchNull = false
	}
	if params.Role != nil && *params.Role != "" {
		roleParam = *params.Role
		roleNull = false
	}
	if params.IsActive != nil {
		isActiveParam = *params.IsActive
		isActiveNull = false
	}

	// For sqlc: empty string / false counts as NULL via the ($1::text IS NULL OR ...) pattern
	// We need to pass the actual values or handle NULL via the pgtype wrappers
	_ = searchNull
	_ = roleNull
	_ = isActiveNull

	listParams := sqlc.ListUsersParams{
		Column1:   searchParam,
		Column2:   roleParam,
		Column3:   isActiveParam,
		Limit:     int32(params.Limit),
		Offset:    int32(offset),
		SortField: params.SortField,
		SortOrder: params.SortOrder,
	}

	users, err := s.queries.ListUsers(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Count with same filters
	countParams := sqlc.CountUsersParams{
		Column1: searchParam,
		Column2: roleParam,
		Column3: isActiveParam,
	}

	total, err := s.queries.CountUsers(ctx, countParams)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Map to response (no password_hash)
	result := &UserListResult{
		Users: make([]UserResponse, len(users)),
		Total: total,
	}
	for i, u := range users {
		result.Users[i] = UserResponse{
			ID:        uuidToString(u.ID),
			Email:     u.Email,
			Role:      u.Role,
			IsActive:  u.IsActive,
			CreatedAt: timestampToString(u.CreatedAt),
			UpdatedAt: timestampToString(u.UpdatedAt),
		}
	}

	return result, nil
}

// GetUser returns a user by ID.
func (s *UserService) GetUser(ctx context.Context, id string) (*UserResponse, error) {
	pgID, err := stringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.queries.GetUserByID(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	resp := UserResponse{
		ID:        uuidToString(user.ID),
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: timestampToString(user.CreatedAt),
		UpdatedAt: timestampToString(user.UpdatedAt),
	}
	return &resp, nil
}

// UpdateUser updates a user.
func (s *UserService) UpdateUser(ctx context.Context, id, email, role string, isActive bool) (*UserResponse, error) {
	pgID, err := stringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:       pgID,
		Column2:  email,
		Column3:  role,
		IsActive: isActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	resp := UserResponse{
		ID:        uuidToString(user.ID),
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: timestampToString(user.CreatedAt),
		UpdatedAt: timestampToString(user.UpdatedAt),
	}
	return &resp, nil
}

// DeleteUser deletes a user.
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	pgID, err := stringToUUID(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = s.queries.DeleteUser(ctx, pgID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func stringToUUID(s string) (pgtype.UUID, error) {
	var pgID pgtype.UUID
	err := pgID.Scan(s)
	return pgID, err
}
