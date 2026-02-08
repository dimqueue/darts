import { Palette } from '../ui/BoxIcon';
import SettingsSection from './SettingsSection';
import { THEME_OPTIONS, type ThemeName } from '../../config/constants';

interface ThemeSelectorProps {
    selectedTheme: ThemeName;
    onThemeChange: (theme: ThemeName) => void;
}

export default function ThemeSelector({ selectedTheme, onThemeChange }: ThemeSelectorProps) {
    return (
        <SettingsSection icon={Palette} title="Color Theme">
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4" role="radiogroup" aria-label="Color theme">
                {THEME_OPTIONS.map((opt) => {
                    const isSelected = selectedTheme === opt.name;
                    return (
                        <button
                            key={opt.name}
                            onClick={() => onThemeChange(opt.name)}
                            className={`p-4 rounded-xl border-2 transition-all duration-200 ${
                                isSelected
                                    ? `${opt.border} ring-2 ring-offset-2 ring-offset-white dark:ring-offset-gray-800 ${opt.ring}`
                                    : 'border-gray-200 dark:border-gray-600 hover:border-gray-300 dark:hover:border-gray-500'
                            }`}
                            role="radio"
                            aria-checked={isSelected}
                        >
                            <div className={`h-12 rounded-lg mb-2 ${opt.preview}`} />
                            <p className="font-medium capitalize text-gray-700 dark:text-gray-200">
                                {opt.name}
                            </p>
                        </button>
                    );
                })}
            </div>
        </SettingsSection>
    );
}
