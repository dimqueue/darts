import { useState, useEffect } from 'react';
import { Trophy, Medal, Crown, Globe } from 'lucide-react';
import api from '../api';
import { useAuth } from '../contexts/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import Layout from '../components/layout/Layout';
import Card from '../components/ui/Card';
import Tabs from '../components/ui/Tabs';
import Pagination from '../components/ui/Pagination';

const TABS = [
    { id: 'global', label: 'Global' },
    { id: 'weekly', label: 'Weekly' },
    { id: 'monthly', label: 'Monthly' },
    { id: 'language', label: 'Language' },
];

const LANGUAGES = [
    { code: 'en', name: 'English' },
    { code: 'ua', name: 'Ukrainian' },
];

const ITEMS_PER_PAGE = 50;

export default function LeaderboardPage() {
    const { user } = useAuth();
    const { theme } = useTheme();
    const [activeTab, setActiveTab] = useState('global');
    const [language, setLanguage] = useState('en');
    const [data, setData] = useState({ users: [], total: 0, current_user_rank: null });
    const [myRank, setMyRank] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [currentPage, setCurrentPage] = useState(1);

    useEffect(() => {
        loadLeaderboard();
        loadMyRank();
    }, [activeTab, language, currentPage]);

    const loadLeaderboard = async () => {
        setLoading(true);
        setError('');
        try {
            const offset = (currentPage - 1) * ITEMS_PER_PAGE;
            const params = { limit: ITEMS_PER_PAGE, offset };

            let response;
            switch (activeTab) {
                case 'weekly':
                    response = await api.getWeeklyLeaderboard(params);
                    break;
                case 'monthly':
                    response = await api.getMonthlyLeaderboard(params);
                    break;
                case 'language':
                    response = await api.getLanguageLeaderboard(language, params);
                    break;
                default:
                    response = await api.getGlobalLeaderboard(params);
            }
            setData(response);
        } catch (err) {
            setError('Failed to load leaderboard: ' + err.message);
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

    const getRankIcon = (rank) => {
        if (rank === 1) return <Crown className="w-5 h-5 text-yellow-500" />;
        if (rank === 2) return <Medal className="w-5 h-5 text-gray-400" />;
        if (rank === 3) return <Medal className="w-5 h-5 text-amber-600" />;
        return <span className="w-5 h-5 text-center text-gray-500">#{rank}</span>;
    };

    const totalPages = Math.ceil(data.total / ITEMS_PER_PAGE);

    return (
        <Layout>
            <div className="max-w-4xl mx-auto space-y-6">
                {/* Header */}
                <Card>
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <Trophy className={`w-8 h-8 ${theme.textColor}`} />
                            <div>
                                <h1 className="text-2xl font-bold text-gray-800">Leaderboard</h1>
                                <p className="text-sm text-gray-500">See how you rank against others</p>
                            </div>
                        </div>
                        {myRank && (
                            <div className="text-right">
                                <p className="text-sm text-gray-500">Your Global Rank</p>
                                <p className={`text-2xl font-bold ${theme.textColor}`}>
                                    #{myRank.global_rank || '-'}
                                </p>
                            </div>
                        )}
                    </div>
                </Card>

                {/* Tabs */}
                <div className="flex items-center gap-4">
                    <Tabs tabs={TABS} activeTab={activeTab} onTabChange={(tab) => {
                        setActiveTab(tab);
                        setCurrentPage(1);
                    }} />

                    {activeTab === 'language' && (
                        <div className="flex items-center gap-2">
                            <Globe className="w-5 h-5 text-gray-400" />
                            <select
                                value={language}
                                onChange={(e) => {
                                    setLanguage(e.target.value);
                                    setCurrentPage(1);
                                }}
                                className={`px-3 py-2 border-2 rounded-xl ${theme.focusBorder} focus:outline-none`}
                            >
                                {LANGUAGES.map((lang) => (
                                    <option key={lang.code} value={lang.code}>
                                        {lang.name}
                                    </option>
                                ))}
                            </select>
                        </div>
                    )}
                </div>

                {/* Leaderboard Table */}
                <Card padding="p-0">
                    {loading ? (
                        <div className="flex items-center justify-center py-20">
                            <div className="animate-spin w-8 h-8 border-4 border-purple-500 border-t-transparent rounded-full" />
                        </div>
                    ) : error ? (
                        <div className="p-6 text-center text-red-500">{error}</div>
                    ) : data.users.length === 0 ? (
                        <div className="p-6 text-center text-gray-500">No data available</div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="w-full">
                                <thead className="bg-gray-50 border-b border-gray-200">
                                    <tr>
                                        <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase">Rank</th>
                                        <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase">Player</th>
                                        <th className="px-6 py-4 text-right text-xs font-semibold text-gray-500 uppercase">Score</th>
                                        <th className="px-6 py-4 text-right text-xs font-semibold text-gray-500 uppercase">Wins</th>
                                        <th className="px-6 py-4 text-right text-xs font-semibold text-gray-500 uppercase">Win Rate</th>
                                        <th className="px-6 py-4 text-right text-xs font-semibold text-gray-500 uppercase">Avg Guesses</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-100">
                                    {data.users.map((entry) => {
                                        const isCurrentUser = entry.username === user?.username;
                                        return (
                                            <tr
                                                key={entry.user_id}
                                                className={`transition-colors ${
                                                    isCurrentUser
                                                        ? `${theme.gradient.replace('bg-', 'bg-opacity-10 bg-')} font-semibold`
                                                        : 'hover:bg-gray-50'
                                                }`}
                                            >
                                                <td className="px-6 py-4">
                                                    <div className="flex items-center gap-2">
                                                        {getRankIcon(entry.rank)}
                                                    </div>
                                                </td>
                                                <td className="px-6 py-4">
                                                    <div className="flex items-center gap-3">
                                                        <div className={`w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold ${theme.gradient}`}>
                                                            {entry.username?.[0]?.toUpperCase() || '?'}
                                                        </div>
                                                        <div>
                                                            <p className="font-medium text-gray-800">
                                                                {entry.name || entry.username}
                                                            </p>
                                                            <p className="text-sm text-gray-500">@{entry.username}</p>
                                                        </div>
                                                    </div>
                                                </td>
                                                <td className="px-6 py-4 text-right font-semibold text-gray-800">
                                                    {entry.total_score?.toLocaleString() || 0}
                                                </td>
                                                <td className="px-6 py-4 text-right text-gray-600">
                                                    {entry.total_wins || 0}
                                                </td>
                                                <td className="px-6 py-4 text-right text-gray-600">
                                                    {entry.win_rate || 0}%
                                                </td>
                                                <td className="px-6 py-4 text-right text-gray-600">
                                                    {entry.average_guesses?.toFixed(1) || '-'}
                                                </td>
                                            </tr>
                                        );
                                    })}
                                </tbody>
                            </table>
                        </div>
                    )}
                </Card>

                {/* Pagination */}
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