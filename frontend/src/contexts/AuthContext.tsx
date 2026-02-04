import { createContext, useContext, useState, type ReactNode } from 'react';
import type { User } from '../types';

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
    return localStorage.getItem('token');
}

function getInitialUser(): User | null {
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
        try {
            return JSON.parse(savedUser);
        } catch {
            localStorage.removeItem('user');
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

    const login = (newToken: string, userData: User) => {
        localStorage.setItem('token', newToken);
        localStorage.setItem('user', JSON.stringify(userData));
        setToken(newToken);
        setUser(userData);
    };

    const logout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        setToken(null);
        setUser(null);
    };

    const updateUser = (userData: Partial<User>) => {
        const updatedUser = { ...user, ...userData } as User;
        localStorage.setItem('user', JSON.stringify(updatedUser));
        setUser(updatedUser);
    };

    const value: AuthContextValue = {
        token,
        user,
        isAuthenticated: !!token,
        loading,
        login,
        logout,
        updateUser,
    };

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
