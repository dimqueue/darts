import { forwardRef, useId } from 'react';
import { useTheme } from '../../contexts/ThemeContext';

interface InputProps {
    label?: string;
    type?: string;
    value?: string;
    onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
    onKeyPress?: (e: React.KeyboardEvent<HTMLInputElement>) => void;
    placeholder?: string;
    disabled?: boolean;
    error?: string;
    className?: string;
    id?: string;
    'aria-label'?: string;
}

const Input = forwardRef<HTMLInputElement, InputProps>(function Input(
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
        id: providedId,
        'aria-label': ariaLabel,
    },
    ref
) {
    const { theme, darkMode } = useTheme();
    const generatedId = useId();
    const inputId = providedId || generatedId;
    const errorId = error ? `${inputId}-error` : undefined;

    return (
        <div className={className}>
            {label && (
                <label
                    htmlFor={inputId}
                    className={`block text-sm font-medium mb-1 ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}
                >
                    {label}
                </label>
            )}
            <input
                ref={ref}
                id={inputId}
                type={type}
                value={value}
                onChange={onChange}
                onKeyPress={onKeyPress}
                placeholder={placeholder}
                disabled={disabled}
                aria-label={ariaLabel}
                aria-invalid={error ? true : undefined}
                aria-describedby={errorId}
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
            {error && (
                <p id={errorId} className="mt-1 text-sm text-red-500" role="alert">
                    {error}
                </p>
            )}
        </div>
    );
});

export default Input;
