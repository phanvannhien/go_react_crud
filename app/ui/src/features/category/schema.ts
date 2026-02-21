import { z } from 'zod';

export const createCategorySchema = z.object({
    name: z.string().trim().min(1, 'Name is required').max(255, 'Name too long'),
});

export const updateCategorySchema = z.object({
    name: z.string().trim().min(1, 'Name is required').max(255, 'Name too long'),
});

export type CreateCategoryFormData = z.infer<typeof createCategorySchema>;
export type UpdateCategoryFormData = z.infer<typeof updateCategorySchema>;
