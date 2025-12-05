import { useState, useEffect, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Trophy, Globe, ChevronUp, ChevronDown } from 'lucide-react';
import api from '@/api';
import { useAuth } from '../contexts/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import Layout from '../components/layout/Layout';
import Card from '../components/ui/Card';
import Pagination from '../components/ui/Pagination';

const LANGUAGES = [
    { code: '', name: 'All Languages' },
    { code: 'en', name: 'English' },
    { code: 'ua', name: 'Ukrainian' },
];

const ITEMS_PER_PAGE = 50;

const getCountryFlag = (countryCode) => {
    if (!countryCode || countryCode.length !== 2) return null;
    const codePoints = countryCode
        .toUpperCase()
        .split('')
        .map(char => 127397 + char.charCodeAt(0));
    return String.fromCodePoint(...codePoints);
};

const RankNavigation = ({ ranks, activeTab, onTabChange, theme }) => {
    const stats = [
        { id: 'global', label: 'Global', value: ranks?.global_rank },
        { id: 'monthly', label: 'Monthly', value: ranks?.monthly_rank },
        { id: 'weekly', label: 'Weekly', value: ranks?.weekly_rank },
        { id: 'daily', label: 'Daily', value: ranks?.daily_rank },
    ];

    return (
        <Card>
            <div className="flex items-center gap-3 mb-4">
                <Trophy className={`w-8 h-8 ${theme.textColor}`} />
                <div>
                    <h1 className="text-2xl font-bold text-gray-800">Leaderboard</h1>
                    <p className="text-sm text-gray-500">Click your rank to switch views</p>
                </div>
            </div>
            <div className="flex gap-3">
                {stats.map((stat) => {
                    const isActive = stat.id === activeTab;
                    return (
                        <button
                            key={stat.id}
                            onClick={() => onTabChange(stat.id)}
                            className={`flex-1 px-4 py-4 rounded-xl text-center border-2 transition-all cursor-pointer ${
                                isActive
                                    ? `${theme.borderColor} ${theme.gradient} text-white shadow-lg scale-105`
                                    : 'border-gray-200 bg-gray-50 hover:border-gray-300 hover:bg-gray-100'
                            }`}
                        >
                            <p className={`text-xs font-medium uppercase tracking-wide ${isActive ? 'text-white/80' : 'text-gray-500'}`}>
                                {stat.label}
                            </p>
                            <p className={`text-2xl font-bold mt-1 ${isActive ? 'text-white' : 'text-gray-800'}`}>
                                {stat.value ? `#${stat.value}` : '-'}
                            </p>
                            {isActive && <div className="w-8 h-1 bg-white/50 rounded-full mx-auto mt-2" />}
                        </button>
                    );
                })}
            </div>
        </Card>
    );
};

const SortableHeader = ({ label, field, sortBy, sortDirection, onSort, align = 'left' }) => {
    const isActive = sortBy === field;
    const alignClass = align === 'right' ? 'text-right justify-end' : 'text-left';

    return (
        <th
            className={`px-6 py-4 text-xs font-semibold text-gray-500 uppercase cursor-pointer hover:bg-gray-100 transition-colors select-none ${alignClass}`}
            onClick={() => onSort(field)}
        >
            <div className={`flex items-center gap-1 ${align === 'right' ? 'justify-end' : ''}`}>
                {label}
                {isActive && (
                    sortDirection === 'asc'
                        ? <ChevronUp className="w-4 h-4" />
                        : <ChevronDown className="w-4 h-4" />
                )}
            </div>
        </th>
    );
};

export default function LeaderboardPage() {
    const { user } = useAuth();
    const { theme } = useTheme();
    const [activeTab, setActiveTab] = useState('global');
    const [language, setLanguage] = useState('');
    const [data, setData] = useState({ users: [], total: 0, current_user_rank: null });
    const [myRank, setMyRank] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const [sortBy, setSortBy] = useState('total_score');
    const [sortDirection, setSortDirection] = useState('desc');

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
            const offset = (currentPage - 1) * ITEMS_PER_PAGE;
            const params = { limit: ITEMS_PER_PAGE, offset };
            const lang = language || null;

            const response = await api.getLeaderboard(activeTab, params, lang);
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

    const handleSort = (field) => {
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
            let aVal = a[sortBy];
            let bVal = b[sortBy];

            if (aVal == null) aVal = 0;
            if (bVal == null) bVal = 0;

            if (sortBy === 'average_guesses') {
                return sortDirection === 'asc' ? aVal - bVal : bVal - aVal;
            }

            return sortDirection === 'asc' ? aVal - bVal : bVal - aVal;
        });

        return sorted;
    }, [data.users, sortBy, sortDirection]);

    const totalPages = Math.ceil(data.total / ITEMS_PER_PAGE);

    const handleTabChange = (tab) => {
        setActiveTab(tab);
        setCurrentPage(1);
    };

    return (
        <Layout>
            <div className="max-w-5xl mx-auto space-y-6">
                {/* Header with Rank Navigation */}
                <RankNavigation
                    ranks={myRank}
                    activeTab={activeTab}
                    onTabChange={handleTabChange}
                    theme={theme}
                />

                {/* Language Filter */}
                <div className="flex items-center justify-end">
                    <div className="flex items-center gap-2">
                        <Globe className="w-5 h-5 text-gray-400" />
                        <select
                            value={language}
                            onChange={(e) => {
                                setLanguage(e.target.value);
                                setCurrentPage(1);
                            }}
                            className={`px-3 py-2 border-2 rounded-xl ${theme.focusBorder} focus:outline-none bg-white`}
                        >
                            {LANGUAGES.map((lang) => (
                                <option key={lang.code} value={lang.code}>
                                    {lang.name}
                                </option>
                            ))}
                        </select>
                    </div>
                </div>

                {/* Leaderboard Table */}
                <Card padding="p-0">
                    {loading ? (
                        <div className="flex items-center justify-center py-20">
                            <div className={`animate-spin w-8 h-8 border-4 ${theme.textColor.replace('text-', 'border-')} border-t-transparent rounded-full`} />
                        </div>
                    ) : error ? (
                        <div className="p-6 text-center text-red-500">{error}</div>
                    ) : sortedUsers.length === 0 ? (
                        <div className="p-6 text-center text-gray-500">No data available</div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="w-full">
                                <thead className="bg-gray-50 border-b border-gray-200">
                                    <tr>
                                        <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase">
                                            Rank
                                        </th>
                                        <th className="px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase">
                                            Player
                                        </th>
                                        <SortableHeader
                                            label="Score"
                                            field="total_score"
                                            sortBy={sortBy}
                                            sortDirection={sortDirection}
                                            onSort={handleSort}
                                            align="right"
                                        />
                                        <SortableHeader
                                            label="Wins"
                                            field="total_wins"
                                            sortBy={sortBy}
                                            sortDirection={sortDirection}
                                            onSort={handleSort}
                                            align="right"
                                        />
                                        <SortableHeader
                                            label="Win Rate"
                                            field="win_rate"
                                            sortBy={sortBy}
                                            sortDirection={sortDirection}
                                            onSort={handleSort}
                                            align="right"
                                        />
                                        <SortableHeader
                                            label="Avg Guesses"
                                            field="average_guesses"
                                            sortBy={sortBy}
                                            sortDirection={sortDirection}
                                            onSort={handleSort}
                                            align="right"
                                        />
                                        <SortableHeader
                                            label="Best Streak"
                                            field="best_win_streak"
                                            sortBy={sortBy}
                                            sortDirection={sortDirection}
                                            onSort={handleSort}
                                            align="right"
                                        />
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-100">
                                    {sortedUsers.map((entry) => {
                                        const isCurrentUser = entry.username === user?.username;
                                        const flag = getCountryFlag(entry.country_code);

                                        return (
                                            <tr
                                                key={entry.user_id}
                                                className={`transition-colors ${
                                                    isCurrentUser
                                                        ? `${theme.gradient.replace('bg-gradient-to-r', 'bg-opacity-20')} bg-purple-50 font-semibold`
                                                        : 'hover:bg-gray-50'
                                                }`}
                                            >
                                                <td className="px-6 py-4">
                                                    <span className="text-gray-600 font-medium">
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
                                                            <p className="font-medium text-gray-800 flex items-center gap-1">
                                                                {flag && <span>{flag}</span>}
                                                                {entry.name || entry.username}
                                                                {isCurrentUser && (
                                                                    <span className={`ml-1 text-xs px-2 py-0.5 rounded-full ${theme.gradient} text-white`}>
                                                                        You
                                                                    </span>
                                                                )}
                                                            </p>
                                                            <p className="text-sm text-gray-500">@{entry.username}</p>
                                                        </div>
                                                    </Link>
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
                                                <td className="px-6 py-4 text-right text-gray-600">
                                                    {entry.best_win_streak || 0}
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
