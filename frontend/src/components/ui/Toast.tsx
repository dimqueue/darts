import { useEffect, useState } from 'react';
import { X, CheckCircle, AlertCircle, Info, AlertTriangle } from './BoxIcon';
import { useToast } from '../../contexts/ToastContext';
import type { ToastItem } from '../../types';

const ICONS = {
    success: CheckCircle,
    error: AlertCircle,
    info: Info,
    warning: AlertTriangle,
} as const;

const STYLES = {
    success: {
        light: 'bg-green-50 border-green-200 text-green-800',
        dark: 'dark:bg-green-900/30 dark:border-green-800 dark:text-green-300',
        icon: 'text-green-500 dark:text-green-400',
    },
    error: {
        light: 'bg-red-50 border-red-200 text-red-800',
        dark: 'dark:bg-red-900/30 dark:border-red-800 dark:text-red-300',
        icon: 'text-red-500 dark:text-red-400',
    },
    info: {
        light: 'bg-blue-50 border-blue-200 text-blue-800',
        dark: 'dark:bg-blue-900/30 dark:border-blue-800 dark:text-blue-300',
        icon: 'text-blue-500 dark:text-blue-400',
    },
    warning: {
        light: 'bg-yellow-50 border-yellow-200 text-yellow-800',
        dark: 'dark:bg-yellow-900/30 dark:border-yellow-800 dark:text-yellow-300',
        icon: 'text-yellow-500 dark:text-yellow-400',
    },
} as const;

interface ToastItemProps {
    toast: ToastItem;
    onRemove: (id: number) => void;
}

function ToastItemComponent({ toast, onRemove }: ToastItemProps) {
    const [isVisible, setIsVisible] = useState(false);
    const [isLeaving, setIsLeaving] = useState(false);

    const Icon = ICONS[toast.type] || Info;
    const style = STYLES[toast.type] || STYLES.info;

    useEffect(() => {
        const timer = setTimeout(() => setIsVisible(true), 10);
        return () => clearTimeout(timer);
    }, []);

    const handleRemove = () => {
        setIsLeaving(true);
        setTimeout(() => onRemove(toast.id), 200);
    };

    return (
        <div
            role="alert"
            aria-live="polite"
            className={`
                flex items-start gap-3 p-4 rounded-xl border shadow-lg backdrop-blur-sm
                transition-all duration-200 ease-out
                ${style.light} ${style.dark}
                ${isVisible && !isLeaving ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'}
            `}
        >
            <Icon className={`w-5 h-5 flex-shrink-0 mt-0.5 ${style.icon}`} />
            <p className="flex-1 text-sm font-medium">{toast.message}</p>
            <button
                onClick={handleRemove}
                className="flex-shrink-0 p-1 rounded-lg hover:bg-black/10 dark:hover:bg-white/10 transition-colors"
                aria-label="Dismiss notification"
            >
                <X className="w-4 h-4" />
            </button>
        </div>
    );
}

export function ToastContainer() {
    const { toasts, removeToast } = useToast();

    if (toasts.length === 0) return null;

    return (
        <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2 max-w-sm w-full pointer-events-none">
            {toasts.map((toast) => (
                <div key={toast.id} className="pointer-events-auto">
                    <ToastItemComponent toast={toast} onRemove={removeToast} />
                </div>
            ))}
        </div>
    );
}

export default ToastContainer;
