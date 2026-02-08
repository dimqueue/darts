import { createContext, useContext, useState, useCallback, useMemo, type ReactNode } from 'react';
import type { User } from '../types';
import { STORAGE_KEYS } from '../config/constants';

interface AuthContextValue {
    token: string | null;
    user: User | null;
    isAuthenticated: boolean;
    loading: boolean;
    login: (newToken: string, userData: User) => void;
    logout: () => void;
    updateUser: (userData: Partial<User>) => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

function getInitialToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.TOKEN);
}

function getInitialUser(): User | null {
    const savedUser = localStorage.getItem(STORAGE_KEYS.USER);
    if (savedUser) {
        try {
            return JSON.parse(savedUser);
        } catch {
            localStorage.removeItem(STORAGE_KEYS.USER);
        }
    }
    return null;
}

interface AuthProviderProps {
    children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
    const [token, setToken] = useState<string | null>(getInitialToken);
    const [user, setUser] = useState<User | null>(getInitialUser);
    const [loading] = useState(false);

    const login = useCallback((newToken: string, userData: User) => {
        localStorage.setItem(STORAGE_KEYS.TOKEN, newToken);
        localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(userData));
        setToken(newToken);
        setUser(userData);
    }, []);

    const logout = useCallback(() => {
        localStorage.removeItem(STORAGE_KEYS.TOKEN);
        localStorage.removeItem(STORAGE_KEYS.USER);
        setToken(null);
        setUser(null);
    }, []);

    const updateUser = useCallback((userData: Partial<User>) => {
        setUser((prev) => {
            const updatedUser = { ...prev, ...userData } as User;
            localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(updatedUser));
            return updatedUser;
        });
    }, []);

    const value = useMemo<AuthContextValue>(
        () => ({
            token,
            user,
            isAuthenticated: !!token,
            loading,
            login,
            logout,
            updateUser,
        }),
        [token, user, loading, login, logout, updateUser]
    );

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}

export default AuthContext;
