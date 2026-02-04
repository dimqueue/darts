import { Globe, ChevronDown } from '../ui/BoxIcon';
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
                <h2 className="font-semibold text-gray-800 dark:text-white">
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
                    className={`w-full px-4 py-3 pr-10 border-2 rounded-xl ${theme.focusBorder} focus:outline-none appearance-none cursor-pointer transition-colors bg-white dark:bg-gray-700 text-gray-800 dark:text-white border-gray-200 dark:border-gray-600 [&>option]:bg-white dark:[&>option]:bg-gray-700`}
                >
                    {LANGUAGES.map((lang) => (
                        <option key={lang.code} value={lang.code}>
                            {lang.name}
                        </option>
                    ))}
                </select>
                <ChevronDown
                    className="absolute right-4 top-1/2 -translate-y-1/2 w-5 h-5 pointer-events-none text-gray-500 dark:text-gray-400"
                    aria-hidden="true"
                />
            </div>
        </Card>
    );
}
