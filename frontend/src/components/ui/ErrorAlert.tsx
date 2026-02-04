import { AlertCircle } from 'lucide-react';
import { useTheme } from '../../contexts/ThemeContext';

interface ErrorAlertProps {
    message: string;
    className?: string;
}

export default function ErrorAlert({ message, className = '' }: ErrorAlertProps) {
    const { darkMode } = useTheme();

    if (!message) return null;

    return (
        <div
            role="alert"
            className={`p-3 rounded-xl text-sm border flex items-start gap-2 ${
                darkMode
                    ? 'bg-red-900/30 border-red-800 text-red-400'
                    : 'bg-red-50 border-red-200 text-red-600'
            } ${className}`}
        >
            <AlertCircle className="w-5 h-5 flex-shrink-0 mt-0.5" />
            <span>{message}</span>
        </div>
    );
}
