import { Globe, ChevronDown } from 'lucide-react';
import Card from '../ui/Card';
import { useTheme } from '../../contexts/ThemeContext';
import { LANGUAGES } from '../../config/constants';

interface LanguageSettingsProps {
    preferredLanguage: string;
    onLanguageChange: (language: string) => void;
}

export default function LanguageSettings({
    preferredLanguage,
    onLanguageChange,
}: LanguageSettingsProps) {
    const { theme, darkMode } = useTheme();

    return (
        <Card>
            <div className="flex items-center gap-3 mb-4">
                <Globe
                    className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`}
                    aria-hidden="true"
                />
                <h2
                    className={`font-semibold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                >
                    Preferred Language
                </h2>
            </div>
            <div className="relative">
                <label htmlFor="preferred-language" className="sr-only">
                    Preferred language
                </label>
                <select
                    id="preferred-language"
                    value={preferredLanguage}
                    onChange={(e) => onLanguageChange(e.target.value)}
                    className={`w-full px-4 py-3 pr-10 border-2 rounded-xl ${theme.focusBorder} focus:outline-none appearance-none cursor-pointer transition-colors ${
                        darkMode
                            ? 'bg-gray-700 text-white border-gray-600 [&>option]:bg-gray-700 [&>option:checked]:bg-violet-600 [&>option:hover]:bg-violet-500'
                            : 'bg-white text-gray-800 border-gray-200 [&>option]:bg-white [&>option:checked]:bg-violet-100 [&>option:hover]:bg-gray-100'
                    }`}
                >
                    {LANGUAGES.map((lang) => (
                        <option key={lang.code} value={lang.code}>
                            {lang.name}
                        </option>
                    ))}
                </select>
                <ChevronDown
                    className={`absolute right-4 top-1/2 -translate-y-1/2 w-5 h-5 pointer-events-none ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                    aria-hidden="true"
                />
            </div>
        </Card>
    );
}
