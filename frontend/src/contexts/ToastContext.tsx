import { createContext, useContext, useState, useCallback, useEffect, useRef, type ReactNode } from 'react';
import type { ToastItem } from '../types';

type ToastType = 'success' | 'error' | 'info' | 'warning';

interface ToastContextValue {
    toasts: ToastItem[];
    addToast: (message: string, type?: ToastType, duration?: number | null) => number;
    removeToast: (id: number) => void;
    success: (message: string, duration?: number) => number;
    error: (message: string, duration?: number) => number;
    info: (message: string, duration?: number) => number;
    warning: (message: string, duration?: number) => number;
}

const ToastContext = createContext<ToastContextValue | null>(null);

let toastId = 0;

interface ToastProviderProps {
    children: ReactNode;
}

export function ToastProvider({ children }: ToastProviderProps) {
    const [toasts, setToasts] = useState<ToastItem[]>([]);
    const timeoutIdsRef = useRef<Map<number, ReturnType<typeof setTimeout>>>(new Map());

    useEffect(() => {
        return () => {
            timeoutIdsRef.current.forEach((timeoutId) => clearTimeout(timeoutId));
            timeoutIdsRef.current.clear();
        };
    }, []);

    const removeToast = useCallback((id: number) => {
        setToasts((prev) => prev.filter((toast) => toast.id !== id));
        const timeoutId = timeoutIdsRef.current.get(id);
        if (timeoutId) {
            clearTimeout(timeoutId);
            timeoutIdsRef.current.delete(id);
        }
    }, []);

    const addToast = useCallback(
        (message: string, type: ToastType = 'info', duration: number | null = null) => {
            const id = ++toastId;

            const defaultDurations: Record<ToastType, number> = {
                success: 3000,
                error: 5000,
                info: 3000,
                warning: 4000,
            };

            const finalDuration = duration ?? defaultDurations[type] ?? 3000;

            const toast: ToastItem = {
                id,
                message,
                type,
                duration: finalDuration,
            };

            setToasts((prev) => [...prev, toast]);

            if (finalDuration > 0) {
                const timeoutId = setTimeout(() => {
                    removeToast(id);
                }, finalDuration);
                timeoutIdsRef.current.set(id, timeoutId);
            }

            return id;
        },
        [removeToast]
    );

    const success = useCallback(
        (message: string, duration?: number) => {
            return addToast(message, 'success', duration);
        },
        [addToast]
    );

    const error = useCallback(
        (message: string, duration?: number) => {
            return addToast(message, 'error', duration);
        },
        [addToast]
    );

    const info = useCallback(
        (message: string, duration?: number) => {
            return addToast(message, 'info', duration);
        },
        [addToast]
    );

    const warning = useCallback(
        (message: string, duration?: number) => {
            return addToast(message, 'warning', duration);
        },
        [addToast]
    );

    const value: ToastContextValue = {
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

export function useToast(): ToastContextValue {
    const context = useContext(ToastContext);
    if (!context) {
        throw new Error('useToast must be used within a ToastProvider');
    }
    return context;
}

export default ToastContext;
