export interface Product {
    id: string;
    name: string;
    description: string;
    price: number;
    category_id: string;
    stock: number;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface ProductListParams {
    limit?: number;
    cursor?: string;
    search?: string;
    category_id?: string;
    is_active?: boolean;
    min_stock?: number;
    min_price?: number;
    max_price?: number;
}

export interface CreateProductRequest {
    name: string;
    description?: string;
    price: number;
    category_id: string;
    stock: number;
    is_active?: boolean;
}

export interface UpdateProductRequest {
    name?: string;
    description?: string;
    price?: number;
    category_id?: string;
    stock?: number;
    is_active?: boolean;
}
