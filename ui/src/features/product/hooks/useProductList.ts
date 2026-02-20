import { useState, useCallback } from 'react';
import { listProductsAPI } from '../api';
import type { Product, ProductListParams } from '../types';

export function useProductList() {
    const [products, setProducts] = useState<Product[]>([]);
    const [nextCursor, setNextCursor] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [filters, setFilters] = useState<Omit<ProductListParams, 'cursor' | 'limit'>>({});

    const fetchProducts = useCallback(async (cursor?: string, append = false) => {
        setLoading(true);
        setError('');
        try {
            const params: ProductListParams = {
                ...filters,
                limit: 20,
                cursor: cursor || undefined,
            };
            const res = await listProductsAPI(params);
            if (append) {
                setProducts(prev => [...prev, ...(res.data || [])]);
            } else {
                setProducts(res.data || []);
            }
            setNextCursor(res.next_cursor || '');
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load products');
        } finally {
            setLoading(false);
        }
    }, [filters]);

    const loadMore = () => {
        if (nextCursor) {
            fetchProducts(nextCursor, true);
        }
    };

    const resetAndFetch = useCallback(() => {
        setProducts([]);
        setNextCursor('');
        fetchProducts();
    }, [fetchProducts]);

    const updateFilters = (newFilters: Partial<Omit<ProductListParams, 'cursor' | 'limit'>>) => {
        setFilters(prev => ({ ...prev, ...newFilters }));
    };

    return {
        products,
        nextCursor,
        loading,
        error,
        filters,
        loadMore,
        updateFilters,
        resetAndFetch,
        fetchProducts,
    };
}
