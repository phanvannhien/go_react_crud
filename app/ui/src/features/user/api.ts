import { apiFetch, buildQueryString, type PaginatedResponse } from '../../lib/api-client';
import type { User, UserListParams, UpdateUserRequest } from './types';

export async function listUsersAPI(params: UserListParams): Promise<PaginatedResponse<User>> {
    return apiFetch<PaginatedResponse<User>>(`/api/users${buildQueryString(params as Record<string, string | number | boolean | undefined>)}`);
}

export async function getUserAPI(id: string): Promise<{ data: User }> {
    return apiFetch<{ data: User }>(`/api/users/${id}`);
}

export async function updateUserAPI(id: string, data: UpdateUserRequest): Promise<{ data: User }> {
    return apiFetch<{ data: User }>(`/api/users/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
    });
}

export async function deleteUserAPI(id: string): Promise<void> {
    await apiFetch(`/api/users/${id}`, { method: 'DELETE' });
}
