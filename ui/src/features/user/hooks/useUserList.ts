import { useState, useEffect, useCallback } from 'react';
import { listUsersAPI } from '../api';
import type { User, UserListParams } from '../types';

export function useUserList() {
    const [users, setUsers] = useState<User[]>([]);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [params, setParams] = useState<UserListParams>({
        page: 1,
        limit: 20,
        sort: 'created_at',
        order: 'desc',
    });

    const fetchUsers = useCallback(async () => {
        setLoading(true);
        setError('');
        try {
            const res = await listUsersAPI(params);
            setUsers(res.data || []);
            setTotal(res.meta?.total || 0);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load users');
        } finally {
            setLoading(false);
        }
    }, [params]);

    useEffect(() => {
        fetchUsers();
    }, [fetchUsers]);

    const setSearch = (search: string) => setParams(p => ({ ...p, search, page: 1 }));
    const setRole = (role: string) => setParams(p => ({ ...p, role: role || undefined, page: 1 }));
    const setIsActive = (val: string) => {
        const isActive = val === '' ? undefined : val === 'true';
        setParams(p => ({ ...p, is_active: isActive, page: 1 }));
    };
    const setPage = (page: number) => setParams(p => ({ ...p, page }));
    const setSort = (sort: string, order: string) => setParams(p => ({ ...p, sort, order }));

    return { users, total, loading, error, params, setSearch, setRole, setIsActive, setPage, setSort, refetch: fetchUsers };
}
