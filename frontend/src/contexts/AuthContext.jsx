import { createContext, useContext, useState } from 'react';

const AuthContext = createContext(null);

function getInitialToken() {
    return localStorage.getItem('token');
}

function getInitialUser() {
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

export function AuthProvider({ children }) {
    const [token, setToken] = useState(getInitialToken);
    const [user, setUser] = useState(getInitialUser);
    const [loading] = useState(false);

    const login = (newToken, userData) => {
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

    const updateUser = (userData) => {
        const updatedUser = { ...user, ...userData };
        localStorage.setItem('user', JSON.stringify(updatedUser));
        setUser(updatedUser);
    };

    const value = {
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

export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}

export default AuthContext;
