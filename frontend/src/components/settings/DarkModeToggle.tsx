import { Moon, Sun } from '../ui/BoxIcon';
import Card from '../ui/Card';
import { useTheme } from '../../contexts/ThemeContext';

interface DarkModeToggleProps {
    onToggle: () => void;
}

export default function DarkModeToggle({ onToggle }: DarkModeToggleProps) {
    const { theme, darkMode } = useTheme();

    return (
        <Card>
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                    {darkMode ? (
                        <Moon
                            className={`w-5 h-5 ${theme.textColorDark}`}
                            aria-hidden="true"
                        />
                    ) : (
                        <Sun className={`w-5 h-5 ${theme.textColor}`} aria-hidden="true" />
                    )}
                    <div>
                        <h2 className="font-semibold text-gray-800 dark:text-white">
                            Dark Mode
                        </h2>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                            {darkMode
                                ? 'Currently using dark theme'
                                : 'Currently using light theme'}
                        </p>
                    </div>
                </div>
                <button
                    onClick={onToggle}
                    className={`relative w-14 h-8 rounded-full transition-colors duration-300 ${
                        darkMode ? theme.gradient : 'bg-gray-300'
                    }`}
                    role="switch"
                    aria-checked={darkMode}
                    aria-label="Toggle dark mode"
                >
                    <span
                        className={`absolute top-1 left-1 w-6 h-6 bg-white rounded-full shadow-md transition-all duration-300 ease-in-out ${
                            darkMode ? 'translate-x-6' : 'translate-x-0'
                        }`}
                    />
                </button>
            </div>
        </Card>
    );
}
