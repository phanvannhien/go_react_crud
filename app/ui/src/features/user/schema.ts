import { z } from 'zod';

export const updateUserSchema = z.object({
    email: z.string().trim().email('Invalid email format').max(255).optional().or(z.literal('')),
    role: z.enum(['user', 'admin']).optional(),
    is_active: z.boolean().optional(),
});

export type UpdateUserFormData = z.infer<typeof updateUserSchema>;
