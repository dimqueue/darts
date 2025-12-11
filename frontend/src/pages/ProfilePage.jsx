import { useState, useEffect } from 'react';
import { User, Trophy, Target, Flame, TrendingUp, Save, Edit2 } from 'lucide-react';
import api from '@/api';
import { useAuth } from '../contexts/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import Layout from '../components/layout/Layout';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';
import StatCard from '../components/ui/StatCard';
import { ProfileSkeleton } from '../components/ui/Skeleton';

export default function ProfilePage() {
    const { user, updateUser } = useAuth();
    const { theme, darkMode } = useTheme();
    const [profile, setProfile] = useState(null);
    const [stats, setStats] = useState(null);
    const [languageStats, setLanguageStats] = useState([]);
    const [loading, setLoading] = useState(true);
    const [editing, setEditing] = useState(false);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    const [editForm, setEditForm] = useState({
        bio: '',
        country_code: '',
    });

    useEffect(() => {
        loadProfile();
    }, []);

    const loadProfile = async () => {
        setLoading(true);
        setError('');
        try {
            const [profileData, statsData, langStatsData] = await Promise.all([
                api.getMyProfile(),
                api.getMyStatistics(),
                api.getMyLanguageStats(),
            ]);
            setProfile(profileData);
            setStats(statsData);
            setLanguageStats(langStatsData?.data || langStatsData || []);
            setEditForm({
                bio: profileData.bio || '',
                country_code: profileData.country_code || '',
            });
        } catch (err) {
            setError('Failed to load profile: ' + err.message);
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async () => {
        setSaving(true);
        setError('');
        try {
            await api.updateMyProfile(editForm);
            setProfile({ ...profile, ...editForm });
            updateUser({ ...user, ...editForm });
            setEditing(false);
        } catch (err) {
            setError('Failed to save: ' + err.message);
        } finally {
            setSaving(false);
        }
    };

    if (loading) {
        return (
            <Layout>
                <ProfileSkeleton />
            </Layout>
        );
    }

    const winRate =
        stats?.total_games > 0 ? ((stats.total_wins / stats.total_games) * 100).toFixed(1) : 0;

    return (
        <Layout>
            <div className="max-w-4xl mx-auto space-y-6">
                {/* Profile Header */}
                <Card>
                    <div className="flex items-start justify-between">
                        <div className="flex items-center gap-4">
                            <div
                                className={`w-20 h-20 ${theme.gradient} rounded-full flex items-center justify-center`}
                            >
                                <User className="w-10 h-10 text-white" />
                            </div>
                            <div>
                                <h1
                                    className={`text-2xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                                >
                                    {profile?.name || user?.username}
                                </h1>
                                <p className={darkMode ? 'text-gray-400' : 'text-gray-500'}>
                                    @{profile?.username || user?.username}
                                </p>
                                {profile?.country_code && (
                                    <p
                                        className={`text-sm mt-1 ${darkMode ? 'text-gray-500' : 'text-gray-400'}`}
                                    >
                                        {profile.country_code}
                                    </p>
                                )}
                            </div>
                        </div>
                        <Button
                            onClick={() => setEditing(!editing)}
                            variant="outline"
                            icon={editing ? null : Edit2}
                        >
                            {editing ? 'Cancel' : 'Edit'}
                        </Button>
                    </div>

                    {editing ? (
                        <div className="mt-6 space-y-4">
                            <Input
                                label="Bio"
                                value={editForm.bio}
                                onChange={(e) => setEditForm({ ...editForm, bio: e.target.value })}
                                placeholder="Tell us about yourself..."
                            />
                            <Input
                                label="Country Code"
                                value={editForm.country_code}
                                onChange={(e) =>
                                    setEditForm({
                                        ...editForm,
                                        country_code: e.target.value.toUpperCase().slice(0, 2),
                                    })
                                }
                                placeholder="US, GB, DE..."
                            />
                            <Button onClick={handleSave} loading={saving} icon={Save}>
                                Save Changes
                            </Button>
                        </div>
                    ) : (
                        profile?.bio && (
                            <p className={`mt-4 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                                {profile.bio}
                            </p>
                        )
                    )}

                    {error && (
                        <div
                            className={`mt-4 p-3 rounded-xl text-sm ${
                                darkMode
                                    ? 'bg-red-900/30 text-red-400 border border-red-800'
                                    : 'bg-red-50 text-red-600'
                            }`}
                        >
                            {error}
                        </div>
                    )}
                </Card>

                {/* Statistics */}
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <StatCard label="Total Games" value={stats?.total_games || 0} icon={Target} />
                    <StatCard label="Wins" value={stats?.total_wins || 0} icon={Trophy} />
                    <StatCard label="Win Rate" value={`${winRate}%`} icon={TrendingUp} />
                    <StatCard
                        label="Best Streak"
                        value={stats?.best_win_streak || 0}
                        icon={Flame}
                    />
                </div>

                {/* More Stats */}
                <Card>
                    <h2
                        className={`text-lg font-semibold mb-4 ${darkMode ? 'text-white' : 'text-gray-800'}`}
                    >
                        Detailed Statistics
                    </h2>
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                        <div>
                            <p
                                className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                            >
                                Current Streak
                            </p>
                            <p
                                className={`text-xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                            >
                                {stats?.current_win_streak || 0}
                            </p>
                        </div>
                        <div>
                            <p
                                className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                            >
                                Avg. Guesses
                            </p>
                            <p
                                className={`text-xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                            >
                                {stats?.average_guesses?.toFixed(1) || '0.0'}
                            </p>
                        </div>
                        <div>
                            <p
                                className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                            >
                                Total Score
                            </p>
                            <p
                                className={`text-xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                            >
                                {stats?.total_score || 0}
                            </p>
                        </div>
                        <div>
                            <p
                                className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                            >
                                Fastest Win
                            </p>
                            <p
                                className={`text-xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                            >
                                {stats?.fastest_win_seconds ? `${stats.fastest_win_seconds}s` : '-'}
                            </p>
                        </div>
                        <div>
                            <p
                                className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                            >
                                Fewest Guesses
                            </p>
                            <p
                                className={`text-xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                            >
                                {stats?.fewest_guesses_win || '-'}
                            </p>
                        </div>
                        <div>
                            <p
                                className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                            >
                                Total Guesses
                            </p>
                            <p
                                className={`text-xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                            >
                                {stats?.total_guesses || 0}
                            </p>
                        </div>
                    </div>
                </Card>

                {/* Language Stats */}
                {languageStats.length > 0 && (
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
                                        <span
                                            className={darkMode ? 'text-gray-400' : 'text-gray-500'}
                                        >
                                            Games:{' '}
                                            <span
                                                className={`font-semibold ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}
                                            >
                                                {lang.games_played}
                                            </span>
                                        </span>
                                        <span
                                            className={darkMode ? 'text-gray-400' : 'text-gray-500'}
                                        >
                                            Won:{' '}
                                            <span
                                                className={`font-semibold ${darkMode ? 'text-green-400' : 'text-green-600'}`}
                                            >
                                                {lang.games_won}
                                            </span>
                                        </span>
                                        <span
                                            className={darkMode ? 'text-gray-400' : 'text-gray-500'}
                                        >
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
                )}
            </div>
        </Layout>
    );
}
