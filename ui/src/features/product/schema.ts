import { z } from 'zod';

export const createProductSchema = z.object({
    name: z.string().trim().min(1, 'Name is required').max(255),
    description: z.string().max(5000).optional().default(''),
    price: z.number().positive('Price must be greater than 0'),
    category_id: z.string().uuid('Invalid category ID'),
    stock: z.number().int().min(0, 'Stock cannot be negative'),
    is_active: z.boolean().optional().default(true),
});

export const updateProductSchema = z.object({
    name: z.string().trim().max(255).optional(),
    description: z.string().max(5000).optional(),
    price: z.number().positive('Price must be greater than 0').optional(),
    category_id: z.string().uuid('Invalid category ID').optional(),
    stock: z.number().int().min(0).optional(),
    is_active: z.boolean().optional(),
});

export type CreateProductFormData = z.infer<typeof createProductSchema>;
export type UpdateProductFormData = z.infer<typeof updateProductSchema>;
