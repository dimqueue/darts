import { memo } from 'react';
import Card from '../ui/Card';
import LoadingSpinner from '../ui/LoadingSpinner';
import SortableHeader from './SortableHeader';
import LeaderboardRow from './LeaderboardRow';
import type { Theme } from '../../types';
import type { LeaderboardEntry } from '../../types/api';

interface LeaderboardTableProps {
    entries: LeaderboardEntry[];
    loading: boolean;
    error: string;
    currentUsername: string | undefined;
    theme: Theme;
    darkMode: boolean;
    sortBy: string;
    sortDirection: 'asc' | 'desc';
    onSort: (field: string) => void;
}

const LeaderboardTable = memo(function LeaderboardTable({
    entries,
    loading,
    error,
    currentUsername,
    theme,
    darkMode,
    sortBy,
    sortDirection,
    onSort,
}: LeaderboardTableProps) {
    if (loading) {
        return (
            <Card padding="p-0">
                <div className="flex items-center justify-center py-20">
                    <LoadingSpinner size="md" />
                </div>
            </Card>
        );
    }

    if (error) {
        return (
            <Card padding="p-0">
                <div className="p-6 text-center text-red-500" role="alert">
                    {error}
                </div>
            </Card>
        );
    }

    if (entries.length === 0) {
        return (
            <Card padding="p-0">
                <div
                    className={`p-6 text-center ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                >
                    No data available
                </div>
            </Card>
        );
    }

    return (
        <Card padding="p-0">
            <div className="overflow-x-auto">
                <table className="w-full">
                    <thead
                        className={`border-b ${darkMode ? 'bg-gray-700 border-gray-600' : 'bg-gray-50 border-gray-200'}`}
                    >
                        <tr>
                            <th
                                className={`px-6 py-4 text-left text-xs font-semibold uppercase ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                            >
                                Rank
                            </th>
                            <th
                                className={`px-6 py-4 text-left text-xs font-semibold uppercase ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                            >
                                Player
                            </th>
                            <SortableHeader
                                label="Score"
                                field="total_score"
                                sortBy={sortBy}
                                sortDirection={sortDirection}
                                onSort={onSort}
                                align="right"
                                darkMode={darkMode}
                            />
                            <SortableHeader
                                label="Wins"
                                field="total_wins"
                                sortBy={sortBy}
                                sortDirection={sortDirection}
                                onSort={onSort}
                                align="right"
                                darkMode={darkMode}
                            />
                            <SortableHeader
                                label="Win Rate"
                                field="win_rate"
                                sortBy={sortBy}
                                sortDirection={sortDirection}
                                onSort={onSort}
                                align="right"
                                darkMode={darkMode}
                            />
                            <SortableHeader
                                label="Avg Guesses"
                                field="average_guesses"
                                sortBy={sortBy}
                                sortDirection={sortDirection}
                                onSort={onSort}
                                align="right"
                                darkMode={darkMode}
                            />
                            <SortableHeader
                                label="Best Streak"
                                field="best_win_streak"
                                sortBy={sortBy}
                                sortDirection={sortDirection}
                                onSort={onSort}
                                align="right"
                                darkMode={darkMode}
                            />
                        </tr>
                    </thead>
                    <tbody
                        className={`divide-y ${darkMode ? 'divide-gray-700' : 'divide-gray-100'}`}
                    >
                        {entries.map((entry) => (
                            <LeaderboardRow
                                key={entry.user_id}
                                entry={entry}
                                isCurrentUser={entry.username === currentUsername}
                                theme={theme}
                                darkMode={darkMode}
                            />
                        ))}
                    </tbody>
                </table>
            </div>
        </Card>
    );
});

export default LeaderboardTable;
