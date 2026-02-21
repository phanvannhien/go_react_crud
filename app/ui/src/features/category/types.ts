export interface Category {
    id: string;
    name: string;
    created_at: string;
    updated_at: string;
}

export interface CategoryListParams {
    page?: number;
    limit?: number;
}

export interface CreateCategoryRequest {
    name: string;
}

export interface UpdateCategoryRequest {
    name: string;
}
