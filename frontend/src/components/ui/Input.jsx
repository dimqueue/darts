import { forwardRef } from 'react';
import { useTheme } from '../../contexts/ThemeContext';

const Input = forwardRef(function Input(
    {
        label,
        type = 'text',
        value,
        onChange,
        onKeyPress,
        placeholder,
        disabled = false,
        error,
        className = '',
    },
    ref
) {
    const { theme, darkMode } = useTheme();

    return (
        <div className={className}>
            {label && (
                <label
                    className={`block text-sm font-medium mb-1 ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}
                >
                    {label}
                </label>
            )}
            <input
                ref={ref}
                type={type}
                value={value}
                onChange={onChange}
                onKeyPress={onKeyPress}
                placeholder={placeholder}
                disabled={disabled}
                className={`w-full px-4 py-3 border-2 rounded-xl focus:outline-none transition-colors
                    ${
                        error
                            ? 'border-red-500 focus:border-red-500'
                            : darkMode
                              ? `border-gray-600 ${theme.focusBorder}`
                              : `border-gray-300 ${theme.focusBorder}`
                    }
                    ${
                        disabled
                            ? darkMode
                                ? 'bg-gray-700 cursor-not-allowed text-gray-400'
                                : 'bg-gray-100 cursor-not-allowed'
                            : darkMode
                              ? 'bg-gray-700 text-white'
                              : 'bg-white text-gray-800'
                    }
                    ${darkMode ? 'placeholder:text-gray-500' : 'placeholder:text-gray-400'}`}
            />
            {error && <p className="mt-1 text-sm text-red-500">{error}</p>}
        </div>
    );
});

export default Input;
