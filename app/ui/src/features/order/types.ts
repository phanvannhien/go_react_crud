export interface OrderItem {
    id: string;
    order_id: string;
    product_id: string;
    quantity: number;
    price: number;
}

export interface Order {
    id: string;
    user_id: string;
    status: 'pending' | 'paid' | 'cancelled' | 'shipped';
    total_amount: number;
    items: OrderItem[];
    created_at: string;
    updated_at: string;
}

export interface OrderListParams {
    page?: number;
    limit?: number;
}

export interface CreateOrderItemRequest {
    product_id: string;
    quantity: number;
    price: number;
}

export interface CreateOrderRequest {
    items: CreateOrderItemRequest[];
}

export interface UpdateOrderStatusRequest {
    status: 'pending' | 'paid' | 'cancelled' | 'shipped';
}
