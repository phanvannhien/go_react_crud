package service

import (
	"context"
	"fmt"

	sqlc "github.com/nhienphan/full_rest_app/db/sqlc"
)

// CategoryService handles category business logic.
type CategoryService struct {
	queries *sqlc.Queries
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(queries *sqlc.Queries) *CategoryService {
	return &CategoryService{queries: queries}
}

// CategoryResponse is the DTO for a category.
type CategoryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CategoryListResult is the paginated response.
type CategoryListResult struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int64              `json:"total"`
}

// CreateCategory creates a new category.
func (s *CategoryService) CreateCategory(ctx context.Context, name string) (*CategoryResponse, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	cat, err := s.queries.CreateCategory(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	resp := toCategoryResponse(cat)
	return &resp, nil
}

// GetCategoryByID retrieves a category by ID.
func (s *CategoryService) GetCategoryByID(ctx context.Context, id string) (*CategoryResponse, error) {
	pgID, err := stringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	cat, err := s.queries.GetCategoryByID(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	resp := toCategoryResponse(cat)
	return &resp, nil
}

// ListCategories retrieves paginated categories.
func (s *CategoryService) ListCategories(ctx context.Context, page, limit int) (*CategoryListResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	cats, err := s.queries.ListCategories(ctx, sqlc.ListCategoriesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	total, err := s.queries.CountCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count categories: %w", err)
	}

	result := &CategoryListResult{
		Categories: make([]CategoryResponse, len(cats)),
		Total:      total,
	}
	for i, c := range cats {
		result.Categories[i] = toCategoryResponse(c)
	}

	return result, nil
}

// UpdateCategory updates a category name.
func (s *CategoryService) UpdateCategory(ctx context.Context, id, name string) (*CategoryResponse, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	pgID, err := stringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	cat, err := s.queries.UpdateCategory(ctx, sqlc.UpdateCategoryParams{
		ID:   pgID,
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	resp := toCategoryResponse(cat)
	return &resp, nil
}

// DeleteCategory deletes a category.
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	pgID, err := stringToUUID(id)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}

	return s.queries.DeleteCategory(ctx, pgID)
}

func toCategoryResponse(c sqlc.Category) CategoryResponse {
	return CategoryResponse{
		ID:        uuidToString(c.ID),
		Name:      c.Name,
		CreatedAt: timestampToString(c.CreatedAt),
		UpdatedAt: timestampToString(c.UpdatedAt),
	}
}
