import { useState, useEffect, useRef, useCallback } from 'react';
import { useProductList } from '../hooks/useProductList';
import { deleteProductAPI, createProductAPI } from '../api';
import { createProductSchema } from '../schema';
import { APIError } from '../../../lib/api-client';

export default function ProductListPage() {
    const { products, nextCursor, loading, error, updateFilters, loadMore, fetchProducts } = useProductList();
    const [searchInput, setSearchInput] = useState('');
    const [showCreateForm, setShowCreateForm] = useState(false);
    const debounceRef = useRef<ReturnType<typeof setTimeout>>(undefined);

    // Initial fetch
    useEffect(() => {
        fetchProducts();
    }, [fetchProducts]);

    // Debounced search
    const handleSearchChange = (value: string) => {
        setSearchInput(value);
        if (debounceRef.current) clearTimeout(debounceRef.current);
        debounceRef.current = setTimeout(() => {
            updateFilters({ search: value || undefined });
        }, 400);
    };



    const handleDelete = async (id: string) => {
        if (!confirm('Delete this product?')) return;
        try {
            await deleteProductAPI(id);
            fetchProducts();
        } catch {
            alert('Failed to delete');
        }
    };

    return (
        <div className="max-w-6xl mx-auto p-6">
            <div className="flex items-center justify-between mb-8">
                <h1 className="text-3xl font-bold bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent">
                    Products
                </h1>
                <button
                    onClick={() => setShowCreateForm(!showCreateForm)}
                    className="px-4 py-2 bg-[var(--color-primary)] text-white rounded-lg text-sm hover:bg-[var(--color-primary-hover)] transition-colors"
                >
                    {showCreateForm ? 'Cancel' : '+ New Product'}
                </button>
            </div>

            {/* Create Form */}
            {showCreateForm && <CreateProductForm onCreated={() => { setShowCreateForm(false); fetchProducts(); }} />}

            {/* Filters */}
            <div className="bg-[var(--color-surface)] rounded-xl p-4 mb-6 border border-[var(--color-border)] flex flex-wrap gap-4 items-end">
                <div className="flex-1 min-w-[200px]">
                    <label htmlFor="product-search" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Search Name</label>
                    <input
                        id="product-search"
                        type="text"
                        value={searchInput}
                        onChange={e => handleSearchChange(e.target.value)}
                        placeholder="Search products..."
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                    />
                </div>

                <div className="min-w-[120px]">
                    <label htmlFor="product-active-filter" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Status</label>
                    <select
                        id="product-active-filter"
                        onChange={e => {
                            const val = e.target.value;
                            updateFilters({ is_active: val === '' ? undefined : val === 'true' });
                            fetchProducts();
                        }}
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                    >
                        <option value="">All</option>
                        <option value="true">Active</option>
                        <option value="false">Inactive</option>
                    </select>
                </div>

                <div className="min-w-[100px]">
                    <label htmlFor="product-min-price" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Min Price</label>
                    <input
                        id="product-min-price"
                        type="number"
                        min="0"
                        step="0.01"
                        onChange={e => {
                            const val = e.target.value;
                            updateFilters({ min_price: val ? parseFloat(val) : undefined });
                        }}
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                    />
                </div>

                <div className="min-w-[100px]">
                    <label htmlFor="product-max-price" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Max Price</label>
                    <input
                        id="product-max-price"
                        type="number"
                        min="0"
                        step="0.01"
                        onChange={e => {
                            const val = e.target.value;
                            updateFilters({ max_price: val ? parseFloat(val) : undefined });
                        }}
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                    />
                </div>
            </div>

            {error && (
                <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg mb-6 text-sm">{error}</div>
            )}

            {/* Product Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {products.map(product => (
                    <div key={product.id} className="bg-[var(--color-surface)] rounded-xl border border-[var(--color-border)] p-5 hover:border-[var(--color-primary)]/50 transition-all group">
                        <div className="flex justify-between items-start mb-3">
                            <h3 className="text-lg font-semibold text-[var(--color-text)] group-hover:text-[var(--color-primary)] transition-colors">
                                {product.name}
                            </h3>
                            <span className={`inline-flex items-center gap-1 text-xs px-2 py-0.5 rounded-full ${product.is_active ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'
                                }`}>
                                {product.is_active ? 'Active' : 'Inactive'}
                            </span>
                        </div>
                        {product.description && (
                            <p className="text-sm text-[var(--color-text-muted)] mb-3 line-clamp-2">{product.description}</p>
                        )}
                        <div className="flex items-center justify-between">
                            <span className="text-2xl font-bold text-emerald-400">${product.price.toFixed(2)}</span>
                            <span className="text-xs text-[var(--color-text-muted)]">Stock: {product.stock}</span>
                        </div>
                        <div className="mt-3 pt-3 border-t border-[var(--color-border)] flex justify-between items-center">
                            <span className="text-xs text-[var(--color-text-muted)]">
                                {new Date(product.created_at).toLocaleDateString()}
                            </span>
                            <button
                                onClick={() => handleDelete(product.id)}
                                className="text-xs text-red-400 hover:text-red-300 transition-colors opacity-0 group-hover:opacity-100"
                            >
                                Delete
                            </button>
                        </div>
                    </div>
                ))}
            </div>

            {loading && <div className="text-center py-8 text-[var(--color-text-muted)]">Loading...</div>}

            {/* Load More (cursor pagination) */}
            {nextCursor && !loading && (
                <div className="text-center mt-8">
                    <button
                        onClick={loadMore}
                        className="px-6 py-2.5 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg text-sm hover:bg-[var(--color-surface-hover)] transition-colors"
                    >
                        Load More
                    </button>
                </div>
            )}

            {!loading && products.length === 0 && (
                <div className="text-center py-16 text-[var(--color-text-muted)]">
                    <p className="text-lg mb-2">No products found</p>
                    <p className="text-sm">Create a product to get started</p>
                </div>
            )}
        </div>
    );
}

// Inline create form component
function CreateProductForm({ onCreated }: { onCreated: () => void }) {
    const [formData, setFormData] = useState({
        name: '', description: '', price: '', category_id: '', stock: '0',
    });
    const [errors, setErrors] = useState<Record<string, string>>({});
    const [serverError, setServerError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleChange = useCallback((e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setFormData(prev => ({ ...prev, [e.target.name]: e.target.value }));
        setErrors(prev => ({ ...prev, [e.target.name]: '' }));
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setServerError('');
        setErrors({});

        const parsed = createProductSchema.safeParse({
            name: formData.name,
            description: formData.description,
            price: parseFloat(formData.price) || 0,
            category_id: formData.category_id,
            stock: parseInt(formData.stock) || 0,
        });

        if (!parsed.success) {
            const fieldErrors: Record<string, string> = {};
            parsed.error.issues.forEach((issue) => {
                if (issue.path[0]) fieldErrors[String(issue.path[0])] = issue.message;
            });
            setErrors(fieldErrors);
            return;
        }

        setLoading(true);
        try {
            await createProductAPI(parsed.data);
            onCreated();
        } catch (err) {
            if (err instanceof APIError) {
                setServerError(err.message);
            } else {
                setServerError('Failed to create product');
            }
        } finally {
            setLoading(false);
        }
    };

    const isValid = formData.name.length > 0 && formData.price.length > 0 && formData.category_id.length > 0;

    return (
        <div className="bg-[var(--color-surface)] rounded-xl border border-[var(--color-border)] p-6 mb-6">
            <h2 className="text-xl font-semibold mb-4">Create Product</h2>
            {serverError && (
                <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg mb-4 text-sm">{serverError}</div>
            )}
            <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label htmlFor="create-name" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Name *</label>
                    <input id="create-name" name="name" value={formData.name} onChange={handleChange}
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]" />
                    {errors.name && <p className="mt-1 text-xs text-red-400">{errors.name}</p>}
                </div>
                <div>
                    <label htmlFor="create-price" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Price *</label>
                    <input id="create-price" name="price" type="number" step="0.01" value={formData.price} onChange={handleChange}
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]" />
                    {errors.price && <p className="mt-1 text-xs text-red-400">{errors.price}</p>}
                </div>
                <div>
                    <label htmlFor="create-category" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Category ID *</label>
                    <input id="create-category" name="category_id" value={formData.category_id} onChange={handleChange} placeholder="UUID"
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]" />
                    {errors.category_id && <p className="mt-1 text-xs text-red-400">{errors.category_id}</p>}
                </div>
                <div>
                    <label htmlFor="create-stock" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Stock</label>
                    <input id="create-stock" name="stock" type="number" value={formData.stock} onChange={handleChange}
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]" />
                    {errors.stock && <p className="mt-1 text-xs text-red-400">{errors.stock}</p>}
                </div>
                <div className="md:col-span-2">
                    <label htmlFor="create-description" className="block text-xs font-medium text-[var(--color-text-muted)] mb-1">Description</label>
                    <textarea id="create-description" name="description" value={formData.description} onChange={handleChange} rows={2}
                        className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] resize-none" />
                </div>
                <div className="md:col-span-2">
                    <button type="submit" disabled={!isValid || loading}
                        className="px-6 py-2.5 bg-[var(--color-primary)] hover:bg-[var(--color-primary-hover)] text-white font-semibold rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed">
                        {loading ? 'Creating...' : 'Create Product'}
                    </button>
                </div>
            </form>
        </div>
    );
}
