import { Palette } from 'lucide-react';
import Card from '../ui/Card';
import { useTheme, THEMES } from '../../contexts/ThemeContext';

interface ThemeSelectorProps {
    selectedTheme: string;
    onThemeChange: (theme: string) => void;
}

const themeColors: Record<string, string> = {
    purple: 'bg-gradient-purple',
    blue: 'bg-gradient-blue',
    green: 'bg-gradient-green',
};

export default function ThemeSelector({ selectedTheme, onThemeChange }: ThemeSelectorProps) {
    const { theme, darkMode } = useTheme();

    return (
        <Card>
            <div className="flex items-center gap-3 mb-4">
                <Palette
                    className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`}
                    aria-hidden="true"
                />
                <h2
                    className={`font-semibold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                >
                    Color Theme
                </h2>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4" role="radiogroup" aria-label="Color theme">
                {Object.keys(THEMES).map((name) => {
                    const isSelected = selectedTheme === name;
                    const selectedThemeObj = THEMES[name];
                    return (
                        <button
                            key={name}
                            onClick={() => onThemeChange(name)}
                            className={`p-4 rounded-xl border-2 transition-all duration-200 ${
                                isSelected
                                    ? `${selectedThemeObj.borderColor} ring-2 ring-offset-2 ${darkMode ? 'ring-offset-gray-800' : 'ring-offset-white'} ${
                                          name === 'purple'
                                              ? 'ring-violet-300'
                                              : name === 'blue'
                                                ? 'ring-blue-300'
                                                : 'ring-emerald-300'
                                      }`
                                    : darkMode
                                      ? 'border-gray-600 hover:border-gray-500'
                                      : 'border-gray-200 hover:border-gray-300'
                            }`}
                            role="radio"
                            aria-checked={isSelected}
                        >
                            <div className={`h-12 rounded-lg mb-2 ${themeColors[name]}`} />
                            <p
                                className={`font-medium capitalize ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}
                            >
                                {name}
                            </p>
                        </button>
                    );
                })}
            </div>
        </Card>
    );
}
