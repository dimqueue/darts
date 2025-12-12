import { useTheme } from '../../contexts/ThemeContext';

export default function Button({
    children,
    onClick,
    disabled = false,
    loading = false,
    variant = 'primary',
    type = 'button',
    className = '',
    icon: Icon,
}) {
    const { theme } = useTheme();

    const baseStyles =
        'px-6 py-3 rounded-xl font-semibold transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed';

    const variants = {
        primary: `${theme.gradient} text-white hover:brightness-110`,
        secondary: `bg-gray-100 text-gray-700 hover:bg-gray-200`,
        outline: `border-2 ${theme.borderColor} ${theme.textColor} bg-transparent ${theme.hoverBg}`,
        danger: 'bg-red-500 text-white hover:bg-red-600',
    };

    return (
        <button
            type={type}
            onClick={onClick}
            disabled={disabled || loading}
            className={`${baseStyles} ${variants[variant]} ${className}`}
        >
            {loading ? (
                <span className="animate-spin w-5 h-5 border-2 border-white border-t-transparent rounded-full" />
            ) : (
                <>
                    {Icon && <Icon className="w-5 h-5" />}
                    {children}
                </>
            )}
        </button>
    );
}
