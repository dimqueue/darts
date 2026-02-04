import { memo } from 'react';
import Card from '../ui/Card';
import LoadingSpinner from '../ui/LoadingSpinner';
import SortableHeader from './SortableHeader';
import LeaderboardRow from './LeaderboardRow';
import type { LeaderboardEntry } from '../../types/api';

interface LeaderboardTableProps {
    entries: LeaderboardEntry[];
    loading: boolean;
    error: string;
    currentUsername: string | undefined;
    sortBy: string;
    sortDirection: 'asc' | 'desc';
    onSort: (field: string) => void;
}

const LeaderboardTable = memo(function LeaderboardTable({
    entries,
    loading,
    error,
    currentUsername,
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
                <div className="p-6 text-center text-gray-500 dark:text-gray-400">
                    No data available
                </div>
            </Card>
        );
    }

    return (
        <Card padding="p-0">
            <div className="overflow-x-auto">
                <table className="w-full">
                    <thead className="border-b bg-gray-50 dark:bg-gray-800 border-gray-200 dark:border-gray-700">
                        <tr>
                            <th className="px-6 py-4 text-left text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">
                                Rank
                            </th>
                            <th className="px-6 py-4 text-left text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">
                                Player
                            </th>
                            <SortableHeader
                                label="Score"
                                field="total_score"
                                sortBy={sortBy}
                                sortDirection={sortDirection}
                                onSort={onSort}
                                align="right"
                            />
                            <SortableHeader
                                label="Wins"
                                field="total_wins"
                                sortBy={sortBy}
                                sortDirection={sortDirection}
                                onSort={onSort}
                                align="right"
                            />
                            <SortableHeader
                                label="Win Rate"
                                field="win_rate"
                                sortBy={sortBy}
                                sortDirection={sortDirection}
                                onSort={onSort}
                                align="right"
                            />
                            <SortableHeader
                                label="Avg Guesses"
                                field="average_guesses"
                                sortBy={sortBy}
                                sortDirection={sortDirection}
                                onSort={onSort}
                                align="right"
                            />
                            <SortableHeader
                                label="Best Streak"
                                field="best_win_streak"
                                sortBy={sortBy}
                                sortDirection={sortDirection}
                                onSort={onSort}
                                align="right"
                            />
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                        {entries.map((entry) => (
                            <LeaderboardRow
                                key={entry.user_id}
                                entry={entry}
                                isCurrentUser={entry.username === currentUsername}
                            />
                        ))}
                    </tbody>
                </table>
            </div>
        </Card>
    );
});

export default LeaderboardTable;
