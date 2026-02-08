import { memo } from 'react';
import { AlertCircle } from './BoxIcon';

interface ErrorAlertProps {
    message: string;
    className?: string;
}

export default memo(function ErrorAlert({ message, className = '' }: ErrorAlertProps) {
    if (!message) return null;

    return (
        <div
            role="alert"
            className={`p-3 rounded-xl text-sm border flex items-start gap-2 bg-red-50 dark:bg-red-900/30 border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 ${className}`}
        >
            <AlertCircle className="w-5 h-5 flex-shrink-0 mt-0.5" />
            <span>{message}</span>
        </div>
    );
});
