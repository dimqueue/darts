import { Moon, Sun } from '../ui/BoxIcon';
import SettingsSection from './SettingsSection';
import { useTheme } from '../../contexts/ThemeContext';

interface DarkModeToggleProps {
    onToggle: () => void;
}

export default function DarkModeToggle({ onToggle }: DarkModeToggleProps) {
    const { darkMode } = useTheme();

    return (
        <SettingsSection icon={darkMode ? Moon : Sun} title="Dark Mode">
            <div className="flex items-center justify-between">
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    {darkMode
                        ? 'Currently using dark theme'
                        : 'Currently using light theme'}
                </p>
                <button
                    onClick={onToggle}
                    className={`relative w-14 h-8 rounded-full transition-colors duration-300 ${
                        darkMode ? 'bg-theme-gradient' : 'bg-gray-300'
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
        </SettingsSection>
    );
}
