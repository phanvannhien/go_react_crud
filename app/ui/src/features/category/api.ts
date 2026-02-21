import { apiFetch, buildQueryString, type PaginatedResponse, type APIResponse } from '../../lib/api-client';
import type { Category, CategoryListParams, CreateCategoryRequest, UpdateCategoryRequest } from './types';

export async function listCategoriesAPI(params: CategoryListParams): Promise<PaginatedResponse<Category>> {
    return apiFetch<PaginatedResponse<Category>>(`/api/categories${buildQueryString(params as Record<string, string | number | boolean | undefined>)}`);
}

export async function getCategoryAPI(id: string): Promise<APIResponse<Category>> {
    return apiFetch<APIResponse<Category>>(`/api/categories/${id}`);
}

export async function createCategoryAPI(data: CreateCategoryRequest): Promise<APIResponse<Category>> {
    return apiFetch<APIResponse<Category>>('/api/categories', {
        method: 'POST',
        body: JSON.stringify(data),
    });
}

export async function updateCategoryAPI(id: string, data: UpdateCategoryRequest): Promise<APIResponse<Category>> {
    return apiFetch<APIResponse<Category>>(`/api/categories/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
    });
}

export async function deleteCategoryAPI(id: string): Promise<void> {
    await apiFetch(`/api/categories/${id}`, { method: 'DELETE' });
}
