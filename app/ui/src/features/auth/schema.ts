import { z } from 'zod';

export const loginSchema = z.object({
    email: z.string().trim().email('Invalid email format').min(1, 'Email is required'),
    password: z.string().min(1, 'Password is required'),
});

export const registerSchema = z.object({
    email: z.string().trim().email('Invalid email format').min(1, 'Email is required').max(255),
    password: z.string().min(8, 'Password must be at least 8 characters').max(128),
    role: z.enum(['user', 'admin']).optional().default('user'),
});

export type LoginFormData = z.infer<typeof loginSchema>;
export type RegisterFormData = z.infer<typeof registerSchema>;
