import type { LucideIcon } from 'lucide-react';
import { useTheme } from '../../contexts/ThemeContext';

interface ButtonProps {
    children?: React.ReactNode;
    onClick?: () => void;
    disabled?: boolean;
    loading?: boolean;
    variant?: 'primary' | 'secondary' | 'outline' | 'danger';
    type?: 'button' | 'submit' | 'reset';
    className?: string;
    icon?: LucideIcon;
    'aria-label'?: string;
}

export default function Button({
    children,
    onClick,
    disabled = false,
    loading = false,
    variant = 'primary',
    type = 'button',
    className = '',
    icon: Icon,
    'aria-label': ariaLabel,
}: ButtonProps) {
    const { theme, darkMode } = useTheme();

    const baseStyles =
        'px-6 py-3 rounded-xl font-semibold transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed';

    const secondaryClass = darkMode
        ? 'bg-gray-700 text-gray-200 hover:bg-gray-600'
        : 'bg-gray-100 text-gray-700 hover:bg-gray-200';

    const outlineClass = darkMode
        ? `border-2 ${theme.borderColor} ${theme.textColorDark} bg-transparent ${theme.hoverBgDark}`
        : `border-2 ${theme.borderColor} ${theme.textColor} bg-transparent ${theme.hoverBg}`;

    const variants = {
        primary: `${theme.gradient} text-white hover:brightness-110`,
        secondary: secondaryClass,
        outline: outlineClass,
        danger: 'bg-red-500 text-white hover:bg-red-600',
    };

    return (
        <button
            type={type}
            onClick={onClick}
            disabled={disabled || loading}
            aria-disabled={disabled || loading}
            aria-busy={loading}
            aria-label={ariaLabel}
            className={`${baseStyles} ${variants[variant]} ${className}`}
        >
            {loading ? (
                <span className="animate-spin w-5 h-5 border-2 border-white border-t-transparent rounded-full" aria-hidden="true" />
            ) : (
                <>
                    {Icon && <Icon className="w-5 h-5" aria-hidden="true" />}
                    {children}
                </>
            )}
        </button>
    );
}
