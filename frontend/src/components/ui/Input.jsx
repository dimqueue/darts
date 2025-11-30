import { useTheme } from '../../contexts/ThemeContext';

export default function Input({
    label,
    type = 'text',
    value,
    onChange,
    onKeyPress,
    placeholder,
    disabled = false,
    error,
    className = '',
}) {
    const { theme } = useTheme();

    return (
        <div className={className}>
            {label && (
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    {label}
                </label>
            )}
            <input
                type={type}
                value={value}
                onChange={onChange}
                onKeyPress={onKeyPress}
                placeholder={placeholder}
                disabled={disabled}
                className={`w-full px-4 py-3 border-2 rounded-xl focus:outline-none transition-colors
                    ${error
                        ? 'border-red-500 focus:border-red-500'
                        : `border-gray-300 ${theme.focusBorder}`
                    }
                    ${disabled ? 'bg-gray-100 cursor-not-allowed' : 'bg-white'}
                    placeholder:text-gray-400`}
            />
            {error && (
                <p className="mt-1 text-sm text-red-500">{error}</p>
            )}
        </div>
    );
}