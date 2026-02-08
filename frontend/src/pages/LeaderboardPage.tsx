import { useState, useEffect, useMemo } from 'react';
import api from '@/api';
import { useAuth } from '../contexts/AuthContext';
import Layout from '../components/layout/Layout';
import Pagination from '../components/ui/Pagination';
import LanguageSelect from '../components/ui/LanguageSelect';
import { RankNavigation, LeaderboardTable } from '../components/leaderboard';
import { LANGUAGES_WITH_ALL, PAGINATION } from '../config/constants';
import type { LeaderboardEntry, RankData } from '../types/api';

export default function LeaderboardPage() {
    const { user } = useAuth();
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
                    <LanguageSelect
                        id="language-filter"
                        value={language}
                        onChange={(val) => {
                            setLanguage(val);
                            setCurrentPage(1);
                        }}
                        languages={LANGUAGES_WITH_ALL}
                    />
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
