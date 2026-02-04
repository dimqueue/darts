import { Link } from 'react-router-dom';
import { memo } from 'react';
import type { Theme } from '../../types';
import type { LeaderboardEntry } from '../../types/api';
import { getCountryFlag } from '../../config/constants';

interface LeaderboardRowProps {
    entry: LeaderboardEntry;
    isCurrentUser: boolean;
    theme: Theme;
    darkMode: boolean;
}

const LeaderboardRow = memo(function LeaderboardRow({
    entry,
    isCurrentUser,
    theme,
    darkMode,
}: LeaderboardRowProps) {
    const flag = getCountryFlag(entry.country_code);

    return (
        <tr
            className={`transition-colors ${
                isCurrentUser
                    ? darkMode
                        ? `${theme.lightBgDark} border-l-4 ${theme.borderColor} font-semibold`
                        : `${theme.lightBg} border-l-4 ${theme.borderColor} font-semibold`
                    : darkMode
                      ? 'hover:bg-gray-700'
                      : 'hover:bg-gray-50'
            }`}
        >
            <td className="px-6 py-4">
                <span className={`font-medium ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                    #{entry.rank}
                </span>
            </td>
            <td className="px-6 py-4">
                <Link
                    to={`/profile/${entry.username}`}
                    className="flex items-center gap-3 hover:opacity-80 transition-opacity"
                >
                    <div
                        className={`w-10 h-10 rounded-full flex items-center justify-center text-white text-sm font-bold ${theme.gradient}`}
                    >
                        {entry.username?.[0]?.toUpperCase() || '?'}
                    </div>
                    <div>
                        <p
                            className={`font-medium flex items-center gap-1 ${darkMode ? 'text-white' : 'text-gray-800'}`}
                        >
                            {flag && <span>{flag}</span>}
                            {entry.name || entry.username}
                        </p>
                        <p className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                            @{entry.username}
                        </p>
                    </div>
                </Link>
            </td>
            <td
                className={`px-6 py-4 text-right font-semibold ${darkMode ? 'text-white' : 'text-gray-800'}`}
            >
                {entry.total_score?.toLocaleString() || 0}
            </td>
            <td className={`px-6 py-4 text-right ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                {entry.total_wins || 0}
            </td>
            <td className={`px-6 py-4 text-right ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                {entry.win_rate || 0}%
            </td>
            <td className={`px-6 py-4 text-right ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                {entry.average_guesses?.toFixed(1) || '-'}
            </td>
            <td className={`px-6 py-4 text-right ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                {entry.best_win_streak || 0}
            </td>
        </tr>
    );
});

export default LeaderboardRow;
