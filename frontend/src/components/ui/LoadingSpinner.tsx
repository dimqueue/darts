import { useTheme } from '../../contexts/ThemeContext';

interface LoadingSpinnerProps {
    size?: 'sm' | 'md' | 'lg';
    className?: string;
}

export default function LoadingSpinner({ size = 'md', className = '' }: LoadingSpinnerProps) {
    const { theme } = useTheme();

    const sizeClasses = {
        sm: 'w-4 h-4 border-2',
        md: 'w-8 h-8 border-4',
        lg: 'w-12 h-12 border-4',
    };

    const borderColor = theme.textColor.replace('text-', 'border-');

    return (
        <div
            className={`animate-spin rounded-full ${sizeClasses[size]} ${borderColor} border-t-transparent ${className}`}
            role="status"
            aria-label="Loading"
        >
            <span className="sr-only">Loading...</span>
        </div>
    );
}
