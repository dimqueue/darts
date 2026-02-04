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
    const { theme } = useTheme();
    const generatedId = useId();
    const inputId = providedId || generatedId;
    const errorId = error ? `${inputId}-error` : undefined;

    return (
        <div className={className}>
            {label && (
                <label
                    htmlFor={inputId}
                    className="block text-sm font-medium mb-1 text-gray-700 dark:text-gray-200"
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
                            : `border-gray-300 dark:border-gray-600 ${theme.focusBorder}`
                    }
                    ${
                        disabled
                            ? 'bg-gray-100 dark:bg-gray-700 cursor-not-allowed text-gray-500 dark:text-gray-400'
                            : 'bg-white dark:bg-gray-700 text-gray-800 dark:text-white'
                    }
                    placeholder:text-gray-400 dark:placeholder:text-gray-500`}
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
