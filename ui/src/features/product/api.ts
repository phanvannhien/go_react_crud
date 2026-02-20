import { apiFetch, buildQueryString, type CursorResponse, type APIResponse } from '../../lib/api-client';
import type { Product, ProductListParams, CreateProductRequest, UpdateProductRequest } from './types';

export async function listProductsAPI(params: ProductListParams): Promise<CursorResponse<Product>> {
    return apiFetch<CursorResponse<Product>>(`/api/products${buildQueryString(params as Record<string, string | number | boolean | undefined>)}`);
}

export async function getProductAPI(id: string): Promise<APIResponse<Product>> {
    return apiFetch<APIResponse<Product>>(`/api/products/${id}`);
}

export async function createProductAPI(data: CreateProductRequest): Promise<APIResponse<Product>> {
    return apiFetch<APIResponse<Product>>('/api/products', {
        method: 'POST',
        body: JSON.stringify(data),
    });
}

export async function updateProductAPI(id: string, data: UpdateProductRequest): Promise<APIResponse<Product>> {
    return apiFetch<APIResponse<Product>>(`/api/products/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
    });
}

export async function deleteProductAPI(id: string): Promise<void> {
    await apiFetch(`/api/products/${id}`, { method: 'DELETE' });
}
