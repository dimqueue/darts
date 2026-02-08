import { memo } from 'react';

interface CardProps {
    children: React.ReactNode;
    className?: string;
    padding?: string;
}

export default memo(function Card({ children, className = '', padding = 'p-6' }: CardProps) {
    return (
        <div
            className={`rounded-card shadow-card ${padding} bg-white dark:bg-gray-800 text-gray-800 dark:text-white ${className}`}
        >
            {children}
        </div>
    );
});
