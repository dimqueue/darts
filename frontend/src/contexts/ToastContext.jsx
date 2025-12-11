import { createContext, useContext, useState, useCallback } from 'react';

const ToastContext = createContext(null);

let toastId = 0;

export function ToastProvider({ children }) {
    const [toasts, setToasts] = useState([]);

    const removeToast = useCallback((id) => {
        setToasts((prev) => prev.filter((toast) => toast.id !== id));
    }, []);

    const addToast = useCallback(
        (message, type = 'info', duration = null) => {
            const id = ++toastId;

            const defaultDurations = {
                success: 3000,
                error: 5000,
                info: 3000,
                warning: 4000,
            };

            const finalDuration = duration ?? defaultDurations[type] ?? 3000;

            const toast = {
                id,
                message,
                type,
                duration: finalDuration,
            };

            setToasts((prev) => [...prev, toast]);

            if (finalDuration > 0) {
                setTimeout(() => {
                    removeToast(id);
                }, finalDuration);
            }

            return id;
        },
        [removeToast]
    );

    const success = useCallback(
        (message, duration) => {
            return addToast(message, 'success', duration);
        },
        [addToast]
    );

    const error = useCallback(
        (message, duration) => {
            return addToast(message, 'error', duration);
        },
        [addToast]
    );

    const info = useCallback(
        (message, duration) => {
            return addToast(message, 'info', duration);
        },
        [addToast]
    );

    const warning = useCallback(
        (message, duration) => {
            return addToast(message, 'warning', duration);
        },
        [addToast]
    );

    const value = {
        toasts,
        addToast,
        removeToast,
        success,
        error,
        info,
        warning,
    };

    return <ToastContext.Provider value={value}>{children}</ToastContext.Provider>;
}

export function useToast() {
    const context = useContext(ToastContext);
    if (!context) {
        throw new Error('useToast must be used within a ToastProvider');
    }
    return context;
}

export default ToastContext;
