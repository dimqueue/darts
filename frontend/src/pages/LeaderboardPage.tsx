import { useState, useEffect, useMemo } from 'react';
import { Globe } from '../components/ui/BoxIcon';
import api from '@/api';
import { useAuth } from '../contexts/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import Layout from '../components/layout/Layout';
import Pagination from '../components/ui/Pagination';
import { RankNavigation, LeaderboardTable } from '../components/leaderboard';
import { LANGUAGES_WITH_ALL, PAGINATION } from '../config/constants';
import type { LeaderboardEntry, RankData } from '../types/api';

export default function LeaderboardPage() {
    const { user } = useAuth();
    const { theme } = useTheme();
    const [activeTab, setActiveTab] = useState('global');
    const [language, setLanguage] = useState('');
    const [data, setData] = useState<{ users: LeaderboardEntry[]; total: number }>({
        users: [],
        total: 0,
    });
    const [myRank, setMyRank] = useState<RankData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const [sortBy, setSortBy] = useState('total_score');
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc');

    useEffect(() => {
        loadLeaderboard();
    }, [activeTab, language, currentPage]);

    useEffect(() => {
        loadMyRank();
    }, []);

    const loadLeaderboard = async () => {
        setLoading(true);
        setError('');
        try {
            const offset = (currentPage - 1) * PAGINATION.DEFAULT_LIMIT;
            const params = { limit: PAGINATION.DEFAULT_LIMIT, offset };
            const lang = language || null;

            const response = await api.getLeaderboard(activeTab, params, lang);
            setData(response);
        } catch (err) {
            setError('Failed to load leaderboard: ' + (err as Error).message);
        } finally {
            setLoading(false);
        }
    };

    const loadMyRank = async () => {
        try {
            const rank = await api.getMyRank();
            setMyRank(rank);
        } catch {
            // Silently fail - rank display is optional
        }
    };

    const handleSort = (field: string) => {
        if (sortBy === field) {
            setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
        } else {
            setSortBy(field);
            setSortDirection('desc');
        }
    };

    const sortedUsers = useMemo(() => {
        if (!data.users || data.users.length === 0) return [];

        const sorted = [...data.users].sort((a, b) => {
            const aVal = (a[sortBy as keyof LeaderboardEntry] as number) ?? 0;
            const bVal = (b[sortBy as keyof LeaderboardEntry] as number) ?? 0;

            return sortDirection === 'asc' ? aVal - bVal : bVal - aVal;
        });

        return sorted;
    }, [data.users, sortBy, sortDirection]);

    const totalPages = Math.ceil(data.total / PAGINATION.DEFAULT_LIMIT);

    const handleTabChange = (tab: string) => {
        setActiveTab(tab);
        setCurrentPage(1);
    };

    return (
        <Layout>
            <div className="space-y-6">
                <RankNavigation
                    ranks={myRank}
                    activeTab={activeTab}
                    onTabChange={handleTabChange}
                />

                <div className="flex items-center justify-end">
                    <div className="flex items-center gap-2">
                        <Globe
                            className="w-5 h-5 text-gray-400 dark:text-gray-500"
                            aria-hidden="true"
                        />
                        <label htmlFor="language-filter" className="sr-only">
                            Filter by language
                        </label>
                        <select
                            id="language-filter"
                            value={language}
                            onChange={(e) => {
                                setLanguage(e.target.value);
                                setCurrentPage(1);
                            }}
                            className={`px-3 py-2 border-2 rounded-xl ${theme.focusBorder} focus:outline-none bg-white dark:bg-gray-800 text-gray-800 dark:text-white border-gray-200 dark:border-gray-600`}
                        >
                            {LANGUAGES_WITH_ALL.map((lang) => (
                                <option key={lang.code} value={lang.code}>
                                    {lang.name}
                                </option>
                            ))}
                        </select>
                    </div>
                </div>

                <LeaderboardTable
                    entries={sortedUsers}
                    loading={loading}
                    error={error}
                    currentUsername={user?.username}
                    sortBy={sortBy}
                    sortDirection={sortDirection}
                    onSort={handleSort}
                />

                {totalPages > 1 && (
                    <Pagination
                        currentPage={currentPage}
                        totalPages={totalPages}
                        onPageChange={setCurrentPage}
                    />
                )}
            </div>
        </Layout>
    );
}
