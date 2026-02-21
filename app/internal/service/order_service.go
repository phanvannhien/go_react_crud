package service

import (
	"context"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhienphan/full_rest_app/db/sqlc"
)

// OrderService handles order business logic.
type OrderService struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

// NewOrderService creates a new OrderService.
func NewOrderService(queries *sqlc.Queries, pool *pgxpool.Pool) *OrderService {
	return &OrderService{queries: queries, pool: pool}
}

// OrderResponse is the order API response.
type OrderResponse struct {
	ID          string              `json:"id"`
	UserID      string              `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

// OrderItemResponse is the order item API response.
type OrderItemResponse struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}

// OrderListResult holds paginated order results.
type OrderListResult struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64           `json:"total"`
}

// OrderCreateParams holds order creation params.
type OrderCreateParams struct {
	UserID string
	Items  []OrderCreateItemParams
}

// OrderCreateItemParams holds order item creation params.
type OrderCreateItemParams struct {
	ProductID string
	Quantity  int32
	Price     float64
}

// CreateOrder creates a new order with items in a transaction.
func (s *OrderService) CreateOrder(ctx context.Context, params OrderCreateParams) (*OrderResponse, error) {
	if len(params.Items) == 0 {
		return nil, fmt.Errorf("order must have at least one item")
	}

	userID, err := stringToUUID(params.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Calculate total amount from items
	var totalAmount float64
	for _, item := range params.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// Begin transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// Create order
	order, err := qtx.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserID:      userID,
		Status:      "pending",
		TotalAmount: orderFloatToNumeric(totalAmount),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	items := make([]OrderItemResponse, len(params.Items))
	for i, item := range params.Items {
		productID, err := stringToUUID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %w", err)
		}

		orderItem, err := qtx.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: productID,
			Quantity:  item.Quantity,
			Price:     orderFloatToNumeric(item.Price),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}

		items[i] = toOrderItemResponse(orderItem)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	resp := toOrderResponse(order)
	resp.Items = items
	return &resp, nil
}

// GetOrderByID returns an order by ID with owner check.
func (s *OrderService) GetOrderByID(ctx context.Context, orderID, userID string) (*OrderResponse, error) {
	pgID, err := stringToUUID(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	order, err := s.queries.GetOrderByID(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Owner check
	if uuidToString(order.UserID) != userID {
		return nil, fmt.Errorf("forbidden: you can only access your own orders")
	}

	// Get items
	orderItems, err := s.queries.GetOrderItemsByOrderID(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	resp := toOrderResponse(order)
	resp.Items = make([]OrderItemResponse, len(orderItems))
	for i, item := range orderItems {
		resp.Items[i] = toOrderItemResponse(item)
	}

	return &resp, nil
}

// UpdateOrderStatus updates an order's status with business rules.
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID, userID, newStatus string) (*OrderResponse, error) {
	pgID, err := stringToUUID(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	// Get current order
	order, err := s.queries.GetOrderByID(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Owner check
	if uuidToString(order.UserID) != userID {
		return nil, fmt.Errorf("forbidden: you can only access your own orders")
	}

	// Business rule: cannot update if shipped
	if order.Status == "shipped" {
		return nil, fmt.Errorf("cannot update order: order has already been shipped")
	}

	// Update status
	updated, err := s.queries.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		ID:     pgID,
		Status: newStatus,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// Get items
	orderItems, err := s.queries.GetOrderItemsByOrderID(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	resp := toOrderResponse(updated)
	resp.Items = make([]OrderItemResponse, len(orderItems))
	for i, item := range orderItems {
		resp.Items[i] = toOrderItemResponse(item)
	}

	return &resp, nil
}

// DeleteOrder deletes an order with business rules.
func (s *OrderService) DeleteOrder(ctx context.Context, orderID, userID string) error {
	pgID, err := stringToUUID(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	// Get current order
	order, err := s.queries.GetOrderByID(ctx, pgID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Owner check
	if uuidToString(order.UserID) != userID {
		return fmt.Errorf("forbidden: you can only access your own orders")
	}

	// Business rule: cannot delete if paid
	if order.Status == "paid" {
		return fmt.Errorf("cannot delete order: paid orders cannot be deleted")
	}

	// Transaction: delete items then order
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	if err := qtx.DeleteOrderItemsByOrderID(ctx, pgID); err != nil {
		return fmt.Errorf("failed to delete order items: %w", err)
	}

	if err := qtx.DeleteOrder(ctx, pgID); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return tx.Commit(ctx)
}

// ListOrders returns a paginated list of orders for a user.
func (s *OrderService) ListOrders(ctx context.Context, userID string, page, limit int) (*OrderListResult, error) {
	// Enforce limit cap
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	pgUserID, err := stringToUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	offset := (page - 1) * limit

	orders, err := s.queries.ListOrdersByUserID(ctx, sqlc.ListOrdersByUserIDParams{
		UserID: pgUserID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	total, err := s.queries.CountOrdersByUserID(ctx, pgUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to count orders: %w", err)
	}

	result := &OrderListResult{
		Orders: make([]OrderResponse, len(orders)),
		Total:  total,
	}

	for i, o := range orders {
		resp := toOrderResponse(o)
		// Get items for each order
		items, err := s.queries.GetOrderItemsByOrderID(ctx, o.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items: %w", err)
		}
		resp.Items = make([]OrderItemResponse, len(items))
		for j, item := range items {
			resp.Items[j] = toOrderItemResponse(item)
		}
		result.Orders[i] = resp
	}

	return result, nil
}

func toOrderResponse(o sqlc.Order) OrderResponse {
	return OrderResponse{
		ID:          uuidToString(o.ID),
		UserID:      uuidToString(o.UserID),
		Status:      o.Status,
		TotalAmount: numericToFloat(o.TotalAmount),
		CreatedAt:   timestampToString(o.CreatedAt),
		UpdatedAt:   timestampToString(o.UpdatedAt),
	}
}

func toOrderItemResponse(item sqlc.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ID:        uuidToString(item.ID),
		OrderID:   uuidToString(item.OrderID),
		ProductID: uuidToString(item.ProductID),
		Quantity:  item.Quantity,
		Price:     numericToFloat(item.Price),
	}
}

func orderFloatToNumeric(f float64) pgtype.Numeric {
	bf := new(big.Float).SetFloat64(f)
	bi, _ := bf.Mul(bf, big.NewFloat(100)).Int(nil)
	return pgtype.Numeric{
		Int:   bi,
		Exp:   -2,
		Valid: true,
	}
}
