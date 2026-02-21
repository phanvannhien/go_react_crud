import { useState, useEffect } from 'react';
import { useOrderList } from '../hooks/useOrderList';
import { createOrderAPI, deleteOrderAPI, updateOrderStatusAPI, getOrderAPI } from '../api';
import { createOrderSchema, updateOrderStatusSchema } from '../schema';
import { APIError } from '../../../lib/api-client';
import type { Order, OrderItem } from '../types';
import SearchableSelect, { type SearchableOption } from '../../../lib/SearchableSelect';
import { listProductsAPI } from '../../product/api';

const STATUS_COLORS: Record<string, string> = {
    pending: 'bg-yellow-500/20 text-yellow-400',
    paid: 'bg-green-500/20 text-green-400',
    cancelled: 'bg-red-500/20 text-red-400',
    shipped: 'bg-blue-500/20 text-blue-400',
};

export default function OrderListPage() {
    const { orders, page, total, totalPages, loading, error, fetchOrders, goToPage, refresh } = useOrderList();
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [detailOrder, setDetailOrder] = useState<Order | null>(null);
    const [detailLoading, setDetailLoading] = useState(false);

    useEffect(() => {
        fetchOrders(1);
    }, []);  // eslint-disable-line react-hooks/exhaustive-deps

    const handleDelete = async (id: string, status: string) => {
        if (status === 'paid') {
            alert('Cannot delete a paid order');
            return;
        }
        if (!confirm('Delete this order?')) return;
        try {
            await deleteOrderAPI(id);
            refresh();
        } catch (err) {
            alert(err instanceof APIError ? err.message : 'Failed to delete');
        }
    };

    const handleStatusUpdate = async (id: string, currentStatus: string, newStatus: string) => {
        if (currentStatus === 'shipped') {
            alert('Cannot update a shipped order');
            return;
        }
        const parsed = updateOrderStatusSchema.safeParse({ status: newStatus });
        if (!parsed.success) return;
        try {
            await updateOrderStatusAPI(id, parsed.data);
            refresh();
        } catch (err) {
            alert(err instanceof APIError ? err.message : 'Failed to update status');
        }
    };

    const handleViewDetail = async (id: string) => {
        setDetailLoading(true);
        try {
            const res = await getOrderAPI(id);
            setDetailOrder(res.data);
        } catch (err) {
            alert(err instanceof APIError ? err.message : 'Failed to load order');
        } finally {
            setDetailLoading(false);
        }
    };

    return (
        <div className="max-w-6xl mx-auto p-6">
            <div className="flex items-center justify-between mb-8">
                <h1 className="text-3xl font-bold bg-gradient-to-r from-amber-400 to-orange-400 bg-clip-text text-transparent">
                    Orders
                </h1>
                <button
                    onClick={() => setShowCreateForm(!showCreateForm)}
                    className="px-4 py-2 bg-[var(--color-primary)] text-white rounded-lg text-sm hover:bg-[var(--color-primary-hover)] transition-colors"
                >
                    {showCreateForm ? 'Cancel' : '+ New Order'}
                </button>
            </div>

            {/* Create Form */}
            {showCreateForm && <CreateOrderForm onCreated={() => { setShowCreateForm(false); refresh(); }} />}

            {/* Error State */}
            {error && (
                <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg mb-6 text-sm">{error}</div>
            )}

            {/* Order Table */}
            {!loading && orders.length > 0 && (
                <div className="bg-[var(--color-surface)] rounded-xl border border-[var(--color-border)] overflow-hidden mb-6">
                    <table className="w-full">
                        <thead>
                            <tr className="border-b border-[var(--color-border)]">
                                <th className="text-left px-5 py-3 text-xs font-medium text-[var(--color-text-muted)] uppercase tracking-wider">Order ID</th>
                                <th className="text-left px-5 py-3 text-xs font-medium text-[var(--color-text-muted)] uppercase tracking-wider">Status</th>
                                <th className="text-right px-5 py-3 text-xs font-medium text-[var(--color-text-muted)] uppercase tracking-wider">Total</th>
                                <th className="text-left px-5 py-3 text-xs font-medium text-[var(--color-text-muted)] uppercase tracking-wider">Date</th>
                                <th className="text-right px-5 py-3 text-xs font-medium text-[var(--color-text-muted)] uppercase tracking-wider">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {orders.map(order => (
                                <tr key={order.id} className="border-b border-[var(--color-border)] last:border-0 hover:bg-[var(--color-surface-hover)] transition-colors">
                                    <td className="px-5 py-4">
                                        <button
                                            onClick={() => handleViewDetail(order.id)}
                                            className="text-sm font-mono text-[var(--color-primary)] hover:underline"
                                        >
                                            {order.id.slice(0, 8)}...
                                        </button>
                                    </td>
                                    <td className="px-5 py-4">
                                        <select
                                            value={order.status}
                                            onChange={e => handleStatusUpdate(order.id, order.status, e.target.value)}
                                            disabled={order.status === 'shipped'}
                                            className={`text-xs px-2.5 py-1 rounded-full border-0 cursor-pointer font-medium ${STATUS_COLORS[order.status] || ''} bg-opacity-20 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]`}
                                        >
                                            <option value="pending">Pending</option>
                                            <option value="paid">Paid</option>
                                            <option value="cancelled">Cancelled</option>
                                            <option value="shipped">Shipped</option>
                                        </select>
                                    </td>
                                    <td className="px-5 py-4 text-right">
                                        <span className="text-lg font-semibold text-emerald-400">${order.total_amount.toFixed(2)}</span>
                                    </td>
                                    <td className="px-5 py-4 text-sm text-[var(--color-text-muted)]">
                                        {new Date(order.created_at).toLocaleDateString()}
                                    </td>
                                    <td className="px-5 py-4 text-right">
                                        <button
                                            onClick={() => handleDelete(order.id, order.status)}
                                            className="text-xs text-red-400 hover:text-red-300 transition-colors"
                                        >
                                            Delete
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}

            {/* Loading State */}
            {loading && <div className="text-center py-8 text-[var(--color-text-muted)]">Loading...</div>}

            {/* Empty State */}
            {!loading && orders.length === 0 && (
                <div className="text-center py-16 text-[var(--color-text-muted)]">
                    <p className="text-lg mb-2">No orders found</p>
                    <p className="text-sm">Create an order to get started</p>
                </div>
            )}

            {/* Pagination */}
            {totalPages > 1 && (
                <div className="flex items-center justify-between mt-6">
                    <span className="text-sm text-[var(--color-text-muted)]">
                        Showing {orders.length} of {total} orders
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

            {/* Order Detail Modal */}
            {(detailOrder || detailLoading) && (
                <OrderDetailModal
                    order={detailOrder}
                    loading={detailLoading}
                    onClose={() => setDetailOrder(null)}
                />
            )}
        </div>
    );
}

// ─── Create Order Form ──────────────────────────────────────────────────────

interface OrderItemRow {
    product_id: string;
    quantity: string;
    price: string;
}

function CreateOrderForm({ onCreated }: { onCreated: () => void }) {
    const [items, setItems] = useState<OrderItemRow[]>([
        { product_id: '', quantity: '1', price: '' },
    ]);
    const [errors, setErrors] = useState<Record<string, string>>({});
    const [serverError, setServerError] = useState('');
    const [loading, setLoading] = useState(false);
    const [productOptions, setProductOptions] = useState<(SearchableOption & { price?: number })[]>([]);
    const [productsLoading, setProductsLoading] = useState(false);

    useEffect(() => {
        setProductsLoading(true);
        listProductsAPI({ limit: 100 }).then(res => {
            setProductOptions((res.data || []).map(p => ({
                value: p.id,
                label: p.name,
                sublabel: `$${p.price.toFixed(2)} · Stock: ${p.stock}`,
                price: p.price,
            })));
        }).catch(() => { }).finally(() => setProductsLoading(false));
    }, []);

    const addItem = () => {
        setItems(prev => [...prev, { product_id: '', quantity: '1', price: '' }]);
    };

    const removeItem = (index: number) => {
        if (items.length <= 1) return;
        setItems(prev => prev.filter((_, i) => i !== index));
    };

    const updateItem = (index: number, field: keyof OrderItemRow, value: string) => {
        setItems(prev => prev.map((item, i) => i === index ? { ...item, [field]: value } : item));
        setErrors(prev => ({ ...prev, [`items.${index}.${field}`]: '' }));
    };

    const calculateTotal = () => {
        return items.reduce((sum, item) => {
            const qty = parseFloat(item.quantity) || 0;
            const price = parseFloat(item.price) || 0;
            return sum + qty * price;
        }, 0);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setServerError('');
        setErrors({});

        const data = {
            items: items.map(item => ({
                product_id: item.product_id.trim(),
                quantity: parseInt(item.quantity) || 0,
                price: parseFloat(item.price) || 0,
            })),
        };

        const parsed = createOrderSchema.safeParse(data);
        if (!parsed.success) {
            const fieldErrors: Record<string, string> = {};
            parsed.error.issues.forEach(issue => {
                const path = issue.path.join('.');
                fieldErrors[path] = issue.message;
            });
            setErrors(fieldErrors);
            return;
        }

        setLoading(true);
        try {
            await createOrderAPI(parsed.data);
            onCreated();
        } catch (err) {
            if (err instanceof APIError) {
                setServerError(err.message);
            } else {
                setServerError('Failed to create order');
            }
        } finally {
            setLoading(false);
        }
    };

    const isValid = items.every(item => item.product_id.trim() && item.quantity && item.price);

    return (
        <div className="bg-[var(--color-surface)] rounded-xl border border-[var(--color-border)] p-6 mb-6">
            <h2 className="text-xl font-semibold mb-4">Create Order</h2>
            {serverError && (
                <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg mb-4 text-sm">{serverError}</div>
            )}
            <form onSubmit={handleSubmit}>
                {/* Items Header */}
                <div className="grid grid-cols-[1fr_100px_100px_40px] gap-3 mb-2">
                    <span className="text-xs font-medium text-[var(--color-text-muted)]">Product *</span>
                    <span className="text-xs font-medium text-[var(--color-text-muted)]">Qty *</span>
                    <span className="text-xs font-medium text-[var(--color-text-muted)]">Price *</span>
                    <span></span>
                </div>

                {/* Item Rows */}
                {items.map((item, index) => (
                    <div key={index} className="grid grid-cols-[1fr_100px_100px_40px] gap-3 mb-3">
                        <div>
                            <SearchableSelect
                                id={`order-item-product-${index}`}
                                options={productOptions}
                                value={item.product_id}
                                onChange={val => {
                                    updateItem(index, 'product_id', val);
                                    // Auto-fill price from selected product
                                    const product = productOptions.find(p => p.value === val);
                                    if (product?.price) {
                                        updateItem(index, 'price', product.price.toString());
                                    }
                                }}
                                placeholder="Select product"
                                loading={productsLoading}
                                error={errors[`items.${index}.product_id`]}
                            />
                        </div>
                        <div>
                            <input
                                id={`order-item-qty-${index}`}
                                type="number"
                                min="1"
                                value={item.quantity}
                                onChange={e => updateItem(index, 'quantity', e.target.value)}
                                className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                            />
                            {errors[`items.${index}.quantity`] && <p className="mt-1 text-xs text-red-400">{errors[`items.${index}.quantity`]}</p>}
                        </div>
                        <div>
                            <input
                                id={`order-item-price-${index}`}
                                type="number"
                                min="0.01"
                                step="0.01"
                                value={item.price}
                                onChange={e => updateItem(index, 'price', e.target.value)}
                                className="w-full px-3 py-2 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-sm text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
                            />
                            {errors[`items.${index}.price`] && <p className="mt-1 text-xs text-red-400">{errors[`items.${index}.price`]}</p>}
                        </div>
                        <button
                            type="button"
                            onClick={() => removeItem(index)}
                            disabled={items.length <= 1}
                            className="px-2 py-2 text-red-400 hover:text-red-300 transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
                        >
                            ✕
                        </button>
                    </div>
                ))}

                {errors['items'] && <p className="mb-3 text-xs text-red-400">{errors['items']}</p>}

                <div className="flex items-center justify-between mt-4">
                    <div className="flex items-center gap-4">
                        <button
                            type="button"
                            onClick={addItem}
                            className="text-sm text-[var(--color-primary)] hover:text-[var(--color-primary-hover)] transition-colors"
                        >
                            + Add Item
                        </button>
                        <span className="text-sm text-[var(--color-text-muted)]">
                            Total: <span className="font-semibold text-emerald-400">${calculateTotal().toFixed(2)}</span>
                        </span>
                    </div>
                    <button
                        type="submit"
                        disabled={!isValid || loading}
                        className="px-6 py-2.5 bg-[var(--color-primary)] hover:bg-[var(--color-primary-hover)] text-white font-semibold rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {loading ? 'Creating...' : 'Create Order'}
                    </button>
                </div>
            </form>
        </div>
    );
}

// ─── Order Detail Modal ─────────────────────────────────────────────────────

function OrderDetailModal({ order, loading, onClose }: { order: Order | null; loading: boolean; onClose: () => void }) {
    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
            <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose}></div>
            <div className="relative bg-[var(--color-surface)] rounded-2xl border border-[var(--color-border)] w-full max-w-lg mx-4 p-6 shadow-2xl">
                <div className="flex items-center justify-between mb-4">
                    <h2 className="text-xl font-semibold">Order Details</h2>
                    <button onClick={onClose} className="text-[var(--color-text-muted)] hover:text-[var(--color-text)] transition-colors text-xl">✕</button>
                </div>

                {loading && <div className="text-center py-8 text-[var(--color-text-muted)]">Loading...</div>}

                {!loading && order && (
                    <>
                        <div className="grid grid-cols-2 gap-4 mb-6">
                            <div>
                                <span className="text-xs text-[var(--color-text-muted)]">Order ID</span>
                                <p className="text-sm font-mono">{order.id}</p>
                            </div>
                            <div>
                                <span className="text-xs text-[var(--color-text-muted)]">Status</span>
                                <p><span className={`text-xs px-2.5 py-1 rounded-full font-medium ${STATUS_COLORS[order.status]}`}>{order.status}</span></p>
                            </div>
                            <div>
                                <span className="text-xs text-[var(--color-text-muted)]">Total Amount</span>
                                <p className="text-lg font-semibold text-emerald-400">${order.total_amount.toFixed(2)}</p>
                            </div>
                            <div>
                                <span className="text-xs text-[var(--color-text-muted)]">Created</span>
                                <p className="text-sm">{new Date(order.created_at).toLocaleString()}</p>
                            </div>
                        </div>

                        <h3 className="text-sm font-semibold text-[var(--color-text-muted)] mb-3 uppercase tracking-wider">Items ({order.items?.length || 0})</h3>
                        <div className="space-y-2">
                            {order.items?.map((item: OrderItem) => (
                                <div key={item.id} className="flex items-center justify-between bg-[var(--color-bg)] rounded-lg px-4 py-3">
                                    <div>
                                        <p className="text-xs font-mono text-[var(--color-text-muted)]">{item.product_id.slice(0, 8)}...</p>
                                        <p className="text-sm">Qty: {item.quantity}</p>
                                    </div>
                                    <span className="font-semibold text-emerald-400">${(item.price * item.quantity).toFixed(2)}</span>
                                </div>
                            ))}
                            {(!order.items || order.items.length === 0) && (
                                <p className="text-sm text-[var(--color-text-muted)] text-center py-4">No items</p>
                            )}
                        </div>
                    </>
                )}
            </div>
        </div>
    );
}
