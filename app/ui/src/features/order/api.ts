import { apiFetch, buildQueryString, type PaginatedResponse, type APIResponse } from '../../lib/api-client';
import type { Order, OrderListParams, CreateOrderRequest, UpdateOrderStatusRequest } from './types';

export async function listOrdersAPI(params: OrderListParams): Promise<PaginatedResponse<Order>> {
    return apiFetch<PaginatedResponse<Order>>(`/api/orders${buildQueryString(params as Record<string, string | number | boolean | undefined>)}`);
}

export async function getOrderAPI(id: string): Promise<APIResponse<Order>> {
    return apiFetch<APIResponse<Order>>(`/api/orders/${id}`);
}

export async function createOrderAPI(data: CreateOrderRequest): Promise<APIResponse<Order>> {
    return apiFetch<APIResponse<Order>>('/api/orders', {
        method: 'POST',
        body: JSON.stringify(data),
    });
}

export async function updateOrderStatusAPI(id: string, data: UpdateOrderStatusRequest): Promise<APIResponse<Order>> {
    return apiFetch<APIResponse<Order>>(`/api/orders/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
    });
}

export async function deleteOrderAPI(id: string): Promise<void> {
    await apiFetch(`/api/orders/${id}`, { method: 'DELETE' });
}
