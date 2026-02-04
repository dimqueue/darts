import { Trophy, Globe } from 'lucide-react';
import Card from '../ui/Card';
import { useTheme } from '../../contexts/ThemeContext';
import { LANGUAGES } from '../../config/constants';

interface GameHeaderProps {
    language: string;
    onLanguageChange: (language: string) => void;
}

export default function GameHeader({ language, onLanguageChange }: GameHeaderProps) {
    const { theme, darkMode } = useTheme();

    return (
        <Card>
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                    <Trophy
                        className={`w-8 h-8 ${darkMode ? theme.textColorDark : theme.textColor}`}
                    />
                    <div>
                        <h1
                            className={`text-2xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                        >
                            Daily Challenge
                        </h1>
                        <p
                            className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                        >
                            Guess the secret word
                        </p>
                    </div>
                </div>
                <div className="flex items-center gap-2">
                    <Globe
                        className={`w-5 h-5 ${darkMode ? 'text-gray-500' : 'text-gray-400'}`}
                        aria-hidden="true"
                    />
                    <label htmlFor="language-select" className="sr-only">
                        Select language
                    </label>
                    <select
                        id="language-select"
                        value={language}
                        onChange={(e) => onLanguageChange(e.target.value)}
                        className={`px-3 py-2 border-2 rounded-xl ${theme.focusBorder} focus:outline-none ${
                            darkMode
                                ? 'bg-gray-700 text-white border-gray-600'
                                : 'bg-white text-gray-800 border-gray-200'
                        }`}
                    >
                        {LANGUAGES.map((lang) => (
                            <option key={lang.code} value={lang.code}>
                                {lang.name}
                            </option>
                        ))}
                    </select>
                </div>
            </div>
        </Card>
    );
}
