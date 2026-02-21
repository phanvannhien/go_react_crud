import { useState, useCallback } from 'react';
import { listOrdersAPI } from '../api';
import type { Order } from '../types';

export function useOrderList() {
    const [orders, setOrders] = useState<Order[]>([]);
    const [page, setPage] = useState(1);
    const [limit] = useState(20);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const fetchOrders = useCallback(async (targetPage?: number) => {
        setLoading(true);
        setError('');
        const currentPage = targetPage ?? page;
        try {
            const res = await listOrdersAPI({ page: currentPage, limit });
            setOrders(res.data || []);
            setTotal(res.meta?.total ?? 0);
            setPage(currentPage);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load orders');
        } finally {
            setLoading(false);
        }
    }, [page, limit]);

    const goToPage = (newPage: number) => {
        fetchOrders(newPage);
    };

    const refresh = () => {
        fetchOrders(page);
    };

    const totalPages = Math.max(1, Math.ceil(total / limit));

    return {
        orders,
        page,
        limit,
        total,
        totalPages,
        loading,
        error,
        fetchOrders,
        goToPage,
        refresh,
    };
}
