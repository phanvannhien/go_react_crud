import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { getToken, clearToken, setToken } from './api-client';

interface AuthUser {
    id: string;
    email: string;
    role: string;
    is_active: boolean;
}

interface AuthContextType {
    user: AuthUser | null;
    token: string | null;
    isAuthenticated: boolean;
    login: (user: AuthUser, token: string) => void;
    logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<AuthUser | null>(null);
    const [token, setTokenState] = useState<string | null>(getToken());

    useEffect(() => {
        const storedToken = getToken();
        if (storedToken) {
            // Decode JWT to get user info (no verification on client)
            try {
                const payload = JSON.parse(atob(storedToken.split('.')[1]));
                setUser({
                    id: payload.user_id,
                    email: payload.email,
                    role: payload.role,
                    is_active: true,
                });
            } catch {
                clearToken();
                setTokenState(null);
            }
        }
    }, []);

    const login = (user: AuthUser, token: string) => {
        setToken(token);
        setTokenState(token);
        setUser(user);
    };

    const logout = () => {
        clearToken();
        setTokenState(null);
        setUser(null);
    };

    return (
        <AuthContext.Provider
            value={{
                user,
                token,
                isAuthenticated: !!token,
                login,
                logout,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}
