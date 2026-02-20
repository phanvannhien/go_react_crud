import { useState } from 'react';
import { useUserList } from '../hooks/useUserList';
import { deleteUserAPI } from '../api';

export default function UserListPage() {
    const { users, total, loading, error, params, setSearch, setRole, setIsActive, setPage, setSort } = useUserList();
    const [searchInput, setSearchInput] = useState('');

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        setSearch(searchInput);
    };

    const handleDelete = async (id: string) => {
        if (!confirm('Are you sure you want to delete this user?')) return;
        try {
            await deleteUserAPI(id);
            window.location.reload();
        } catch {
            alert('Failed to delete user');
        }
    };

    const totalPages = Math.ceil(total / (params.limit || 20));
    const currentPage = params.page || 1;

    const sortFields = [
        { value: 'created_at', label: 'Created' },
        { value: 'email', label: 'Email' },
        { value: 'role', label: 'Role' },
    ];

    return (
        <div className="max-w-6xl mx-auto p-6">
            <div className="flex items-center justify-between mb-8">
                <h1 className="text-3xl font-bold bg-gradient-to-r from-indigo-400 to-cyan-400 bg-clip-text text-transparent">
                    Users
                </h1>
                <span className="text-sm text-[var(--color-text-muted)]">{total} total</span>
            </div>

            {/* Filters bar */}
            <div className="bg-[var(--color-surface)] rounded-xl p-4 mb-6 border border-[var(--color-border)] flex flex-wrap gap-4 items-end">
                <form onSubmit={handleSearch} className="flex-1 min-w-[200px]">
                    <label htmlFor="user-search" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Search Email</label>
                    <div className="flex gap-2">
                        <input
                            id="user-search"
                            type="text"
                            value={searchInput}
                            onChange={e => setSearchInput(e.target.value)}
                            placeholder="Search by email..."
                            className="flex-1 px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                        />
                        <button type="submit" className="px-4 py-2 bg-[var(--color-primary)] text-white rounded-lg text-sm hover:bg-[var(--color-primary-hover)] transition-colors">
                            Search
                        </button>
                    </div>
                </form>

                <div className="min-w-[120px]">
                    <label htmlFor="user-role-filter" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Role</label>
                    <select
                        id="user-role-filter"
                        value={params.role || ''}
                        onChange={e => setRole(e.target.value)}
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                    >
                        <option value="">All</option>
                        <option value="user">User</option>
                        <option value="admin">Admin</option>
                    </select>
                </div>

                <div className="min-w-[120px]">
                    <label htmlFor="user-active-filter" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Status</label>
                    <select
                        id="user-active-filter"
                        value={params.is_active === undefined ? '' : String(params.is_active)}
                        onChange={e => setIsActive(e.target.value)}
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                    >
                        <option value="">All</option>
                        <option value="true">Active</option>
                        <option value="false">Inactive</option>
                    </select>
                </div>

                <div className="min-w-[140px]">
                    <label htmlFor="user-sort" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Sort</label>
                    <select
                        id="user-sort"
                        value={`${params.sort}:${params.order}`}
                        onChange={e => {
                            const [sort, order] = e.target.value.split(':');
                            setSort(sort, order);
                        }}
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                    >
                        {sortFields.map(f => (
                            <optgroup key={f.value} label={f.label}>
                                <option value={`${f.value}:asc`}>{f.label} ↑</option>
                                <option value={`${f.value}:desc`}>{f.label} ↓</option>
                            </optgroup>
                        ))}
                    </select>
                </div>
            </div>

            {/* Error */}
            {error && (
                <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg mb-6 text-sm">{error}</div>
            )}

            {/* Table */}
            <div className="bg-[var(--color-surface)] rounded-xl border border-[var(--color-border)] overflow-hidden">
                <table className="w-full">
                    <thead>
                        <tr className="border-b border-[var(--color-border)]">
                            <th className="text-left px-4 py-3 text-xs font-semibold text-[var(--color-text-muted)] uppercase tracking-wider">Email</th>
                            <th className="text-left px-4 py-3 text-xs font-semibold text-[var(--color-text-muted)] uppercase tracking-wider">Role</th>
                            <th className="text-left px-4 py-3 text-xs font-semibold text-[var(--color-text-muted)] uppercase tracking-wider">Status</th>
                            <th className="text-left px-4 py-3 text-xs font-semibold text-[var(--color-text-muted)] uppercase tracking-wider">Created</th>
                            <th className="text-right px-4 py-3 text-xs font-semibold text-[var(--color-text-muted)] uppercase tracking-wider">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {loading ? (
                            <tr><td colSpan={5} className="text-center py-12 text-[var(--color-text-muted)]">Loading...</td></tr>
                        ) : users.length === 0 ? (
                            <tr><td colSpan={5} className="text-center py-12 text-[var(--color-text-muted)]">No users found</td></tr>
                        ) : (
                            users.map(user => (
                                <tr key={user.id} className="border-b border-[var(--color-border)] hover:bg-[var(--color-surface-hover)] transition-colors">
                                    <td className="px-4 py-3 text-sm font-medium">{user.email}</td>
                                    <td className="px-4 py-3">
                                        <span className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium ${user.role === 'admin' ? 'bg-purple-500/20 text-purple-400' : 'bg-blue-500/20 text-blue-400'
                                            }`}>
                                            {user.role}
                                        </span>
                                    </td>
                                    <td className="px-4 py-3">
                                        <span className={`inline-flex items-center gap-1 text-xs ${user.is_active ? 'text-green-400' : 'text-red-400'}`}>
                                            <span className={`w-1.5 h-1.5 rounded-full ${user.is_active ? 'bg-green-400' : 'bg-red-400'}`}></span>
                                            {user.is_active ? 'Active' : 'Inactive'}
                                        </span>
                                    </td>
                                    <td className="px-4 py-3 text-sm text-[var(--color-text-muted)]">
                                        {new Date(user.created_at).toLocaleDateString()}
                                    </td>
                                    <td className="px-4 py-3 text-right">
                                        <button
                                            onClick={() => handleDelete(user.id)}
                                            className="text-xs text-red-400 hover:text-red-300 transition-colors"
                                        >
                                            Delete
                                        </button>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
                <div className="flex items-center justify-between mt-6">
                    <span className="text-sm text-[var(--color-text-muted)]">
                        Page {currentPage} of {totalPages}
                    </span>
                    <div className="flex gap-2">
                        <button
                            onClick={() => setPage(currentPage - 1)}
                            disabled={currentPage <= 1}
                            className="px-3 py-1.5 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-sm disabled:opacity-50 hover:bg-[var(--color-surface-hover)] transition-colors"
                        >
                            Previous
                        </button>
                        <button
                            onClick={() => setPage(currentPage + 1)}
                            disabled={currentPage >= totalPages}
                            className="px-3 py-1.5 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-sm disabled:opacity-50 hover:bg-[var(--color-surface-hover)] transition-colors"
                        >
                            Next
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}
