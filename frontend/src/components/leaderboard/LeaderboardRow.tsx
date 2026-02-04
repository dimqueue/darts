import { Link } from 'react-router-dom';
import { memo } from 'react';
import { useTheme } from '../../contexts/ThemeContext';
import type { LeaderboardEntry } from '../../types/api';
import { getCountryFlag } from '../../config/constants';

interface LeaderboardRowProps {
    entry: LeaderboardEntry;
    isCurrentUser: boolean;
}

const LeaderboardRow = memo(function LeaderboardRow({
    entry,
    isCurrentUser,
}: LeaderboardRowProps) {
    const { theme, darkMode } = useTheme();
    const flag = getCountryFlag(entry.country_code);

    // Keep ternary for isCurrentUser - uses theme colors
    const rowClass = isCurrentUser
        ? darkMode
            ? `${theme.lightBgDark} border-l-4 ${theme.borderColor} font-semibold`
            : `${theme.lightBg} border-l-4 ${theme.borderColor} font-semibold`
        : 'hover:bg-gray-50 dark:hover:bg-gray-700';

    return (
        <tr className={rowClass}>
            <td className="px-6 py-4">
                <span className="font-medium text-gray-600 dark:text-gray-300">
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
                        <p className="font-medium flex items-center gap-1 text-gray-800 dark:text-white">
                            {flag && <span>{flag}</span>}
                            {entry.name || entry.username}
                        </p>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                            @{entry.username}
                        </p>
                    </div>
                </Link>
            </td>
            <td className="px-6 py-4 text-right font-semibold text-gray-800 dark:text-white">
                {entry.total_score?.toLocaleString() || 0}
            </td>
            <td className="px-6 py-4 text-right text-gray-600 dark:text-gray-300">
                {entry.total_wins || 0}
            </td>
            <td className="px-6 py-4 text-right text-gray-600 dark:text-gray-300">
                {entry.win_rate || 0}%
            </td>
            <td className="px-6 py-4 text-right text-gray-600 dark:text-gray-300">
                {entry.average_guesses?.toFixed(1) || '-'}
            </td>
            <td className="px-6 py-4 text-right text-gray-600 dark:text-gray-300">
                {entry.best_win_streak || 0}
            </td>
        </tr>
    );
});

export default LeaderboardRow;
