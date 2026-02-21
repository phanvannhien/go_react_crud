package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nhienphan/full_rest_app/db/sqlc"
)

// ProductService handles product business logic.
type ProductService struct {
	queries *sqlc.Queries
}

// NewProductService creates a new ProductService.
func NewProductService(queries *sqlc.Queries) *ProductService {
	return &ProductService{queries: queries}
}

// ProductResponse is the product API response.
type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`
	Stock       int32   `json:"stock"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ProductListResult holds cursor-paginated product results.
type ProductListResult struct {
	Products   []ProductResponse `json:"products"`
	NextCursor string            `json:"next_cursor"`
}

// ProductListParams holds list query params.
type ProductListParams struct {
	Search     *string
	CategoryID *string
	IsActive   *bool
	MinStock   *int32
	MinPrice   *float64
	MaxPrice   *float64
	Cursor     *string // RFC3339 timestamp cursor
	Limit      int
}

// ProductCreateParams holds product creation params.
type ProductCreateParams struct {
	Name        string
	Description string
	Price       float64
	CategoryID  string
	Stock       int32
	IsActive    bool
}

// ProductUpdateParams holds product update params.
type ProductUpdateParams struct {
	ID          string
	Name        string
	Description string
	Price       float64
	CategoryID  string
	Stock       int32
	IsActive    bool
}

// ListProducts returns a cursor-paginated list of products.
func (s *ProductService) ListProducts(ctx context.Context, params ProductListParams) (*ProductListResult, error) {
	// Enforce limit cap
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 20
	}

	// Build list params
	listParams := sqlc.ListProductsParams{
		Limit: int32(params.Limit),
	}

	// Search filter
	if params.Search != nil && *params.Search != "" {
		listParams.Column1 = *params.Search
	}

	// Category filter
	if params.CategoryID != nil && *params.CategoryID != "" {
		catID, err := stringToUUID(*params.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %w", err)
		}
		listParams.Column2 = catID
	}

	// IsActive filter
	if params.IsActive != nil {
		if *params.IsActive {
			listParams.Column3 = "true"
		} else {
			listParams.Column3 = "false"
		}
	}

	// MinStock filter
	if params.MinStock != nil {
		listParams.Column4 = *params.MinStock
	}

	// MinPrice filter
	if params.MinPrice != nil {
		listParams.Column5 = floatToNumeric(*params.MinPrice)
	}

	// MaxPrice filter
	if params.MaxPrice != nil {
		listParams.Column6 = floatToNumeric(*params.MaxPrice)
	}

	// Cursor filter
	if params.Cursor != nil && *params.Cursor != "" {
		t, err := time.Parse(time.RFC3339, *params.Cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor format, must be RFC3339: %w", err)
		}
		listParams.Column7 = pgtype.Timestamptz{Time: t, Valid: true}
	}

	products, err := s.queries.ListProducts(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	result := &ProductListResult{
		Products: make([]ProductResponse, len(products)),
	}

	for i, p := range products {
		result.Products[i] = toProductResponse(p)
	}

	// Set next cursor from last item's created_at
	if len(products) == int(params.Limit) && len(products) > 0 {
		last := products[len(products)-1]
		if last.CreatedAt.Valid {
			result.NextCursor = last.CreatedAt.Time.Format(time.RFC3339Nano)
		}
	}

	return result, nil
}

// GetProduct returns a product by ID.
func (s *ProductService) GetProduct(ctx context.Context, id string) (*ProductResponse, error) {
	pgID, err := stringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	product, err := s.queries.GetProductByID(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	resp := toProductResponse(product)
	return &resp, nil
}

// CreateProduct creates a new product.
func (s *ProductService) CreateProduct(ctx context.Context, params ProductCreateParams) (*ProductResponse, error) {
	catID, err := stringToUUID(params.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	product, err := s.queries.CreateProduct(ctx, sqlc.CreateProductParams{
		Name:        params.Name,
		Description: pgtype.Text{String: params.Description, Valid: params.Description != ""},
		Price:       floatToNumeric(params.Price),
		CategoryID:  catID,
		Stock:       params.Stock,
		IsActive:    params.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	resp := toProductResponse(product)
	return &resp, nil
}

// UpdateProduct updates a product.
func (s *ProductService) UpdateProduct(ctx context.Context, params ProductUpdateParams) (*ProductResponse, error) {
	pgID, err := stringToUUID(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	catID, err := stringToUUID(params.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	product, err := s.queries.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:          pgID,
		Column2:     params.Name,
		Description: pgtype.Text{String: params.Description, Valid: params.Description != ""},
		Price:       floatToNumeric(params.Price),
		CategoryID:  catID,
		Stock:       params.Stock,
		IsActive:    params.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	resp := toProductResponse(product)
	return &resp, nil
}

// DeleteProduct deletes a product.
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	pgID, err := stringToUUID(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	return s.queries.DeleteProduct(ctx, pgID)
}

func toProductResponse(p sqlc.Product) ProductResponse {
	return ProductResponse{
		ID:          uuidToString(p.ID),
		Name:        p.Name,
		Description: p.Description.String,
		Price:       numericToFloat(p.Price),
		CategoryID:  uuidToString(p.CategoryID),
		Stock:       p.Stock,
		IsActive:    p.IsActive,
		CreatedAt:   timestampToString(p.CreatedAt),
		UpdatedAt:   timestampToString(p.UpdatedAt),
	}
}

func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}

func floatToNumeric(f float64) pgtype.Numeric {
	// Convert float64 to pgtype.Numeric
	bf := new(big.Float).SetFloat64(f)
	// Use 2 decimal precision for price
	bi, _ := bf.Mul(bf, big.NewFloat(100)).Int(nil)
	return pgtype.Numeric{
		Int:   bi,
		Exp:   -2,
		Valid: true,
	}
}
