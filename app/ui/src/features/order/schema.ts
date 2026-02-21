import { z } from 'zod';

export const createOrderItemSchema = z.object({
    product_id: z.string().uuid('Invalid product ID'),
    quantity: z.number().int().positive('Quantity must be greater than 0'),
    price: z.number().positive('Price must be greater than 0'),
});

export const createOrderSchema = z.object({
    items: z.array(createOrderItemSchema).min(1, 'At least one item is required'),
});

export const updateOrderStatusSchema = z.object({
    status: z.enum(['pending', 'paid', 'cancelled', 'shipped'], {
        message: 'Status must be pending, paid, cancelled, or shipped',
    }),
});

export type CreateOrderFormData = z.infer<typeof createOrderSchema>;
export type UpdateOrderStatusFormData = z.infer<typeof updateOrderStatusSchema>;
