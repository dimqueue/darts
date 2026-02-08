import { Globe, ChevronDown } from '../ui/BoxIcon';
import SettingsSection from './SettingsSection';
import { LANGUAGES } from '../../config/constants';

interface LanguageSettingsProps {
    preferredLanguage: string;
    onLanguageChange: (language: string) => void;
}

export default function LanguageSettings({
    preferredLanguage,
    onLanguageChange,
}: LanguageSettingsProps) {
    return (
        <SettingsSection icon={Globe} title="Preferred Language">
            <div className="relative">
                <label htmlFor="preferred-language" className="sr-only">
                    Preferred language
                </label>
                <select
                    id="preferred-language"
                    value={preferredLanguage}
                    onChange={(e) => onLanguageChange(e.target.value)}
                    className="w-full px-4 py-3 pr-10 border-2 rounded-xl focus:border-theme-border focus:outline-none appearance-none cursor-pointer transition-colors bg-white dark:bg-gray-700 text-gray-800 dark:text-white border-gray-200 dark:border-gray-600 [&>option]:bg-white dark:[&>option]:bg-gray-700"
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
        </SettingsSection>
    );
}
