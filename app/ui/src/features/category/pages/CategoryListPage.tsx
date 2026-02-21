import { useState, useEffect } from 'react';
import { useCategoryList } from '../hooks/useCategoryList';
import { createCategoryAPI, updateCategoryAPI, deleteCategoryAPI } from '../api';
import { createCategorySchema, updateCategorySchema } from '../schema';
import { APIError } from '../../../lib/api-client';

export default function CategoryListPage() {
    const { categories, page, total, totalPages, loading, error, fetchCategories, goToPage, refresh } = useCategoryList();
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [editingId, setEditingId] = useState<string | null>(null);
    const [editName, setEditName] = useState('');
    const [editError, setEditError] = useState('');
    const [editLoading, setEditLoading] = useState(false);

    useEffect(() => {
        fetchCategories(1);
    }, []);  // eslint-disable-line react-hooks/exhaustive-deps

    const handleDelete = async (id: string) => {
        if (!confirm('Delete this category? Products using it may be affected.')) return;
        try {
            await deleteCategoryAPI(id);
            refresh();
        } catch (err) {
            alert(err instanceof APIError ? err.message : 'Failed to delete');
        }
    };

    const startEdit = (id: string, currentName: string) => {
        setEditingId(id);
        setEditName(currentName);
        setEditError('');
    };

    const cancelEdit = () => {
        setEditingId(null);
        setEditName('');
        setEditError('');
    };

    const handleUpdate = async (id: string) => {
        setEditError('');
        const parsed = updateCategorySchema.safeParse({ name: editName });
        if (!parsed.success) {
            setEditError(parsed.error.issues[0]?.message || 'Invalid name');
            return;
        }

        setEditLoading(true);
        try {
            await updateCategoryAPI(id, parsed.data);
            cancelEdit();
            refresh();
        } catch (err) {
            setEditError(err instanceof APIError ? err.message : 'Failed to update');
        } finally {
            setEditLoading(false);
        }
    };

    return (
        <div className="max-w-4xl mx-auto p-6">
            <div className="flex items-center justify-between mb-8">
                <h1 className="text-3xl font-bold bg-gradient-to-r from-violet-400 to-fuchsia-400 bg-clip-text text-transparent">
                    Categories
                </h1>
                <button
                    onClick={() => setShowCreateForm(!showCreateForm)}
                    className="px-4 py-2 bg-[var(--color-primary)] text-white rounded-lg text-sm hover:bg-[var(--color-primary-hover)] transition-colors"
                >
                    {showCreateForm ? 'Cancel' : '+ New Category'}
                </button>
            </div>

            {/* Create Form */}
            {showCreateForm && <CreateCategoryForm onCreated={() => { setShowCreateForm(false); refresh(); }} />}

            {/* Error State */}
            {error && (
                <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg mb-6 text-sm">{error}</div>
            )}

            {/* Category List */}
            {!loading && categories.length > 0 && (
                <div className="bg-[var(--color-surface)] rounded-xl border border-[var(--color-border)] overflow-hidden mb-6">
                    {categories.map((cat, index) => (
                        <div
                            key={cat.id}
                            className={`flex items-center justify-between px-5 py-4 hover:bg-[var(--color-surface-hover)] transition-colors ${index < categories.length - 1 ? 'border-b border-[var(--color-border)]' : ''
                                }`}
                        >
                            {editingId === cat.id ? (
                                /* Editing Mode */
                                <div className="flex-1 flex items-center gap-3">
                                    <input
                                        value={editName}
                                        onChange={e => { setEditName(e.target.value); setEditError(''); }}
                                        onKeyDown={e => { if (e.key === 'Enter') handleUpdate(cat.id); if (e.key === 'Escape') cancelEdit(); }}
                                        autoFocus
                                        className="flex-1 px-3 py-1.5 bg-[var(--color-bg)] border border-[var(--color-primary)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                                    />
                                    <button
                                        onClick={() => handleUpdate(cat.id)}
                                        disabled={editLoading}
                                        className="px-3 py-1.5 bg-[var(--color-primary)] text-white rounded-lg text-xs hover:bg-[var(--color-primary-hover)] transition-colors disabled:opacity-50"
                                    >
                                        {editLoading ? '...' : 'Save'}
                                    </button>
                                    <button
                                        onClick={cancelEdit}
                                        className="px-3 py-1.5 text-[var(--color-text-muted)] text-xs hover:text-[var(--color-text)] transition-colors"
                                    >
                                        Cancel
                                    </button>
                                    {editError && <span className="text-xs text-red-400">{editError}</span>}
                                </div>
                            ) : (
                                /* Display Mode */
                                <>
                                    <div className="flex items-center gap-4">
                                        <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-violet-500/20 to-fuchsia-500/20 flex items-center justify-center">
                                            <span className="text-sm font-semibold text-violet-400">{cat.name.charAt(0).toUpperCase()}</span>
                                        </div>
                                        <div>
                                            <p className="text-sm font-medium text-[var(--color-text)]">{cat.name}</p>
                                            <p className="text-xs text-[var(--color-text-muted)]">
                                                Created {new Date(cat.created_at).toLocaleDateString()}
                                            </p>
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-3">
                                        <button
                                            onClick={() => startEdit(cat.id, cat.name)}
                                            className="text-xs text-[var(--color-primary)] hover:text-[var(--color-primary-hover)] transition-colors"
                                        >
                                            Edit
                                        </button>
                                        <button
                                            onClick={() => handleDelete(cat.id)}
                                            className="text-xs text-red-400 hover:text-red-300 transition-colors"
                                        >
                                            Delete
                                        </button>
                                    </div>
                                </>
                            )}
                        </div>
                    ))}
                </div>
            )}

            {/* Loading State */}
            {loading && <div className="text-center py-8 text-[var(--color-text-muted)]">Loading...</div>}

            {/* Empty State */}
            {!loading && categories.length === 0 && (
                <div className="text-center py-16 text-[var(--color-text-muted)]">
                    <p className="text-lg mb-2">No categories found</p>
                    <p className="text-sm">Create a category to organize your products</p>
                </div>
            )}

            {/* Pagination */}
            {totalPages > 1 && (
                <div className="flex items-center justify-between mt-6">
                    <span className="text-sm text-[var(--color-text-muted)]">
                        Showing {categories.length} of {total} categories
                    </span>
                    <div className="flex items-center gap-2">
                        <button
                            onClick={() => goToPage(page - 1)}
                            disabled={page <= 1}
                            className="px-3 py-1.5 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-sm hover:bg-[var(--color-surface-hover)] transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                        >
                            Prev
                        </button>
                        {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                            let pageNum: number;
                            if (totalPages <= 5) {
                                pageNum = i + 1;
                            } else if (page <= 3) {
                                pageNum = i + 1;
                            } else if (page >= totalPages - 2) {
                                pageNum = totalPages - 4 + i;
                            } else {
                                pageNum = page - 2 + i;
                            }
                            return (
                                <button
                                    key={pageNum}
                                    onClick={() => goToPage(pageNum)}
                                    className={`px-3 py-1.5 rounded-lg text-sm transition-colors ${pageNum === page
                                        ? 'bg-[var(--color-primary)] text-white'
                                        : 'bg-[var(--color-surface)] border border-[var(--color-border)] hover:bg-[var(--color-surface-hover)]'
                                        }`}
                                >
                                    {pageNum}
                                </button>
                            );
                        })}
                        <button
                            onClick={() => goToPage(page + 1)}
                            disabled={page >= totalPages}
                            className="px-3 py-1.5 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-sm hover:bg-[var(--color-surface-hover)] transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                        >
                            Next
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}

// ─── Create Category Form ───────────────────────────────────────────────────

function CreateCategoryForm({ onCreated }: { onCreated: () => void }) {
    const [name, setName] = useState('');
    const [error, setError] = useState('');
    const [serverError, setServerError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setServerError('');
        setError('');

        const parsed = createCategorySchema.safeParse({ name });
        if (!parsed.success) {
            setError(parsed.error.issues[0]?.message || 'Invalid name');
            return;
        }

        setLoading(true);
        try {
            await createCategoryAPI(parsed.data);
            onCreated();
        } catch (err) {
            if (err instanceof APIError) {
                setServerError(err.message);
            } else {
                setServerError('Failed to create category');
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="bg-[var(--color-surface)] rounded-xl border border-[var(--color-border)] p-6 mb-6">
            <h2 className="text-xl font-semibold mb-4">Create Category</h2>
            {serverError && (
                <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg mb-4 text-sm">{serverError}</div>
            )}
            <form onSubmit={handleSubmit} className="flex items-end gap-4">
                <div className="flex-1">
                    <label htmlFor="create-category-name" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Name *</label>
                    <input
                        id="create-category-name"
                        value={name}
                        onChange={e => { setName(e.target.value); setError(''); }}
                        placeholder="Category name"
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                    />
                    {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
                </div>
                <button
                    type="submit"
                    disabled={!name.trim() || loading}
                    className="px-6 py-2.5 bg-[var(--color-primary)] hover:bg-[var(--color-primary-hover)] text-white font-semibold rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    {loading ? 'Creating...' : 'Create'}
                </button>
            </form>
        </div>
    );
}
