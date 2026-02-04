import Card from '../ui/Card';
import { useTheme } from '../../contexts/ThemeContext';
import type { LanguageStatsData } from '../../types/api';

interface LanguageStatsProps {
    languageStats: LanguageStatsData[];
}

export default function LanguageStats({ languageStats }: LanguageStatsProps) {
    const { darkMode } = useTheme();

    if (languageStats.length === 0) return null;

    return (
        <Card>
            <h2
                className={`text-lg font-semibold mb-4 ${darkMode ? 'text-white' : 'text-gray-800'}`}
            >
                Stats by Language
            </h2>
            <div className="space-y-3">
                {languageStats.map((lang) => (
                    <div
                        key={lang.language}
                        className={`flex items-center justify-between p-3 rounded-xl ${
                            darkMode ? 'bg-gray-700/50' : 'bg-gray-50'
                        }`}
                    >
                        <div className="flex items-center gap-3">
                            <span
                                className={`font-medium uppercase ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}
                            >
                                {lang.language}
                            </span>
                        </div>
                        <div className="flex gap-6 text-sm">
                            <span className={darkMode ? 'text-gray-400' : 'text-gray-500'}>
                                Games:{' '}
                                <span
                                    className={`font-semibold ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}
                                >
                                    {lang.games_played}
                                </span>
                            </span>
                            <span className={darkMode ? 'text-gray-400' : 'text-gray-500'}>
                                Won:{' '}
                                <span
                                    className={`font-semibold ${darkMode ? 'text-green-400' : 'text-green-600'}`}
                                >
                                    {lang.games_won}
                                </span>
                            </span>
                            <span className={darkMode ? 'text-gray-400' : 'text-gray-500'}>
                                Avg:{' '}
                                <span
                                    className={`font-semibold ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}
                                >
                                    {lang.average_guesses?.toFixed(1)}
                                </span>
                            </span>
                        </div>
                    </div>
                ))}
            </div>
        </Card>
    );
}
