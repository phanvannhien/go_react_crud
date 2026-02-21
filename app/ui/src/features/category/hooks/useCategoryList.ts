import { useState, useCallback } from 'react';
import { listCategoriesAPI } from '../api';
import type { Category } from '../types';

export function useCategoryList() {
    const [categories, setCategories] = useState<Category[]>([]);
    const [page, setPage] = useState(1);
    const [limit] = useState(20);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const fetchCategories = useCallback(async (targetPage?: number) => {
        setLoading(true);
        setError('');
        const currentPage = targetPage ?? page;
        try {
            const res = await listCategoriesAPI({ page: currentPage, limit });
            setCategories(res.data || []);
            setTotal(res.meta?.total ?? 0);
            setPage(currentPage);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load categories');
        } finally {
            setLoading(false);
        }
    }, [page, limit]);

    const goToPage = (newPage: number) => {
        fetchCategories(newPage);
    };

    const refresh = () => {
        fetchCategories(page);
    };

    const totalPages = Math.max(1, Math.ceil(total / limit));

    return {
        categories,
        page,
        limit,
        total,
        totalPages,
        loading,
        error,
        fetchCategories,
        goToPage,
        refresh,
    };
}
