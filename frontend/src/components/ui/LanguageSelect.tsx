import { Globe } from './BoxIcon';
import type { Language } from '../../types';

interface LanguageSelectProps {
    value: string;
    onChange: (value: string) => void;
    languages: Language[];
    id: string;
}

export default function LanguageSelect({ value, onChange, languages, id }: LanguageSelectProps) {
    return (
        <div className="flex items-center gap-2">
            <Globe
                className="w-5 h-5 text-gray-400 dark:text-gray-500"
                aria-hidden="true"
            />
            <label htmlFor={id} className="sr-only">
                Select language
            </label>
            <select
                id={id}
                value={value}
                onChange={(e) => onChange(e.target.value)}
                className="px-3 py-2 border-2 rounded-xl focus:border-theme-border focus:outline-none bg-white dark:bg-gray-700 text-gray-800 dark:text-white border-gray-200 dark:border-gray-600"
            >
                {languages.map((lang) => (
                    <option key={lang.code} value={lang.code}>
                        {lang.name}
                    </option>
                ))}
            </select>
        </div>
    );
}
