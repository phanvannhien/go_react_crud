import { apiFetch, type APIResponse } from '../../lib/api-client';
import type { AuthResponse, LoginRequest, RegisterRequest } from './types';

export async function loginAPI(data: LoginRequest): Promise<AuthResponse> {
    const res = await apiFetch<APIResponse<AuthResponse>>('/api/auth/login', {
        method: 'POST',
        body: JSON.stringify(data),
    });
    return res.data;
}

export async function registerAPI(data: RegisterRequest): Promise<AuthResponse> {
    const res = await apiFetch<APIResponse<AuthResponse>>('/api/auth/register', {
        method: 'POST',
        body: JSON.stringify(data),
    });
    return res.data;
}
