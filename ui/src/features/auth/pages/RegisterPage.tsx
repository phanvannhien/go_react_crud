import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { registerSchema, type RegisterFormData } from '../schema';
import { registerAPI } from '../api';
import { useAuth } from '../../../lib/auth-context';
import { APIError } from '../../../lib/api-client';

export default function RegisterPage() {
    const [formData, setFormData] = useState<RegisterFormData>({ email: '', password: '', role: 'user' });
    const [errors, setErrors] = useState<Record<string, string>>({});
    const [serverError, setServerError] = useState('');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();
    const { login } = useAuth();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
        setErrors({ ...errors, [e.target.name]: '' });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setServerError('');
        setErrors({});

        const result = registerSchema.safeParse(formData);
        if (!result.success) {
            const fieldErrors: Record<string, string> = {};
            result.error.issues.forEach((issue) => {
                if (issue.path[0]) fieldErrors[String(issue.path[0])] = issue.message;
            });
            setErrors(fieldErrors);
            return;
        }

        setLoading(true);
        try {
            const res = await registerAPI(result.data);
            login(res.user, res.token);
            navigate('/');
        } catch (err) {
            if (err instanceof APIError) {
                if (err.details) {
                    setErrors(err.details);
                } else {
                    setServerError(err.message);
                }
            } else {
                setServerError('An unexpected error occurred');
            }
        } finally {
            setLoading(false);
        }
    };

    const isValid = formData.email.length > 0 && formData.password.length >= 8;

    return (
        <div className="min-h-screen flex items-center justify-center px-4">
            <div className="w-full max-w-md">
                <div className="bg-[var(--color-surface)] rounded-2xl p-8 shadow-2xl border border-[var(--color-border)]">
                    <h1 className="text-3xl font-bold text-center mb-2 bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent">
                        Create Account
                    </h1>
                    <p className="text-center text-[var(--color-text-muted)] mb-8">Join us today</p>

                    {serverError && (
                        <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg mb-6 text-sm">
                            {serverError}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-5">
                        <div>
                            <label htmlFor="register-email" className="block text-sm font-medium text-[var(--color-text-muted)] mb-1.5">Email</label>
                            <input
                                id="register-email"
                                name="email"
                                type="email"
                                value={formData.email}
                                onChange={handleChange}
                                className="w-full px-4 py-3 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] transition-all"
                                placeholder="you@example.com"
                            />
                            {errors.email && <p className="mt-1 text-sm text-red-400">{errors.email}</p>}
                        </div>

                        <div>
                            <label htmlFor="register-password" className="block text-sm font-medium text-[var(--color-text-muted)] mb-1.5">Password</label>
                            <input
                                id="register-password"
                                name="password"
                                type="password"
                                value={formData.password}
                                onChange={handleChange}
                                className="w-full px-4 py-3 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] transition-all"
                                placeholder="Minimum 8 characters"
                            />
                            {errors.password && <p className="mt-1 text-sm text-red-400">{errors.password}</p>}
                        </div>

                        <div>
                            <label htmlFor="register-role" className="block text-sm font-medium text-[var(--color-text-muted)] mb-1.5">Role</label>
                            <select
                                id="register-role"
                                name="role"
                                value={formData.role}
                                onChange={handleChange}
                                className="w-full px-4 py-3 bg-[var(--color-bg)] border border-[var(--color-border)] rounded-lg text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] transition-all"
                            >
                                <option value="user">User</option>
                                <option value="admin">Admin</option>
                            </select>
                        </div>

                        <button
                            type="submit"
                            disabled={!isValid || loading}
                            className="w-full py-3 bg-[var(--color-primary)] hover:bg-[var(--color-primary-hover)] text-white font-semibold rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {loading ? 'Creating account...' : 'Create Account'}
                        </button>
                    </form>

                    <p className="mt-6 text-center text-sm text-[var(--color-text-muted)]">
                        Already have an account?{' '}
                        <Link to="/login" className="text-[var(--color-primary)] hover:underline font-medium">
                            Sign in
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    );
}
