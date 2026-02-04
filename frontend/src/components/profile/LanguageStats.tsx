import Card from '../ui/Card';
import type { LanguageStatsData } from '../../types/api';

interface LanguageStatsProps {
    languageStats: LanguageStatsData[];
}

export default function LanguageStats({ languageStats }: LanguageStatsProps) {
    if (languageStats.length === 0) return null;

    return (
        <Card>
            <h2 className="text-lg font-semibold mb-4 text-gray-800 dark:text-white">
                Stats by Language
            </h2>
            <div className="space-y-3">
                {languageStats.map((lang) => (
                    <div
                        key={lang.language}
                        className="flex items-center justify-between p-3 rounded-xl bg-gray-50 dark:bg-gray-700/50"
                    >
                        <div className="flex items-center gap-3">
                            <span className="font-medium uppercase text-gray-700 dark:text-gray-200">
                                {lang.language}
                            </span>
                        </div>
                        <div className="flex gap-6 text-sm">
                            <span className="text-gray-500 dark:text-gray-400">
                                Games:{' '}
                                <span className="font-semibold text-gray-700 dark:text-gray-200">
                                    {lang.games_played}
                                </span>
                            </span>
                            <span className="text-gray-500 dark:text-gray-400">
                                Won:{' '}
                                <span className="font-semibold text-green-600 dark:text-green-400">
                                    {lang.games_won}
                                </span>
                            </span>
                            <span className="text-gray-500 dark:text-gray-400">
                                Avg:{' '}
                                <span className="font-semibold text-gray-700 dark:text-gray-200">
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
