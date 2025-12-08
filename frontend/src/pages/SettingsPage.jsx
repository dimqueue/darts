import { useState, useEffect } from 'react';
import { Settings, Palette, Globe, Bell, Eye, Save, Moon, Sun, ChevronDown } from 'lucide-react';
import api from '@/api';
import { useTheme, THEMES } from '../contexts/ThemeContext';
import Layout from '../components/layout/Layout';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import { SettingsSkeleton } from '../components/ui/Skeleton';

const LANGUAGES = [
    { code: 'en', name: 'English' },
    { code: 'ua', name: 'Ukrainian' },
];

export default function SettingsPage() {
    const { themeName, setTheme, theme, darkMode, setDarkMode } = useTheme();
    const [settings, setSettings] = useState({
        preferred_language: 'en',
        theme: themeName,
        sound_enabled: true,
        email_notifications: true,
        show_profile_public: true,
        show_stats_public: true,
    });
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    useEffect(() => {
        setSettings(prev => ({ ...prev, theme: themeName }));
    }, [themeName]);

    useEffect(() => {
        loadSettings();
    }, []);

    const loadSettings = async () => {
        setLoading(true);
        try {
            const data = await api.getMySettings();
            setSettings(prev => ({ ...prev, ...data, theme: themeName }));
        } catch (err) {
            setError('Failed to load settings: ' + err.message);
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async () => {
        setSaving(true);
        setError('');
        setSuccess('');
        try {
            await api.updateMySettings(settings);
            setTheme(settings.theme);
            setSuccess('Settings saved successfully!');
            setTimeout(() => setSuccess(''), 3000);
        } catch (err) {
            setError('Failed to save settings: ' + err.message);
        } finally {
            setSaving(false);
        }
    };

    const handleChange = (key, value) => {
        setSettings({ ...settings, [key]: value });
        if (key === 'theme') {
            setTheme(value);
        }
    };

    const themeColors = {
        purple: 'bg-gradient-purple',
        blue: 'bg-gradient-blue',
        green: 'bg-gradient-green',
    };

    if (loading) {
        return (
            <Layout>
                <SettingsSkeleton />
            </Layout>
        );
    }

    return (
        <Layout>
            <div className="max-w-2xl mx-auto space-y-6">
                {/* Header */}
                <Card>
                    <div className="flex items-center gap-3">
                        <Settings className={`w-8 h-8 ${darkMode ? theme.textColorDark : theme.textColor}`} />
                        <div>
                            <h1 className={`text-2xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}>Settings</h1>
                            <p className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>Customize your experience</p>
                        </div>
                    </div>
                </Card>

                {/* Dark Mode */}
                <Card>
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            {darkMode ? (
                                <Moon className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`} />
                            ) : (
                                <Sun className={`w-5 h-5 ${theme.textColor}`} />
                            )}
                            <div>
                                <h2 className="font-semibold">Dark Mode</h2>
                                <p className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                                    {darkMode ? 'Currently using dark theme' : 'Currently using light theme'}
                                </p>
                            </div>
                        </div>
                        <button
                            onClick={() => setDarkMode(!darkMode)}
                            className={`relative w-14 h-8 rounded-full transition-colors duration-300 ${
                                darkMode ? theme.gradient : 'bg-gray-300'
                            }`}
                            role="switch"
                            aria-checked={darkMode}
                        >
                            <span
                                className={`absolute top-1 left-1 w-6 h-6 bg-white rounded-full shadow-md transition-all duration-300 ease-in-out ${
                                    darkMode ? 'translate-x-6' : 'translate-x-0'
                                }`}
                            />
                        </button>
                    </div>
                </Card>

                {/* Theme Selection */}
                <Card>
                    <div className="flex items-center gap-3 mb-4">
                        <Palette className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`} />
                        <h2 className={`font-semibold ${darkMode ? 'text-white' : 'text-gray-800'}`}>Color Theme</h2>
                    </div>
                    <div className="grid grid-cols-3 gap-4">
                        {Object.keys(THEMES).map((name) => {
                            const isSelected = settings.theme === name;
                            const selectedTheme = THEMES[name];
                            return (
                                <button
                                    key={name}
                                    onClick={() => handleChange('theme', name)}
                                    className={`p-4 rounded-xl border-2 transition-all duration-200 ${
                                        isSelected
                                            ? `${selectedTheme.borderColor} ring-2 ring-offset-2 ${darkMode ? 'ring-offset-gray-800' : 'ring-offset-white'} ${
                                                name === 'purple' ? 'ring-violet-300' :
                                                name === 'blue' ? 'ring-blue-300' :
                                                'ring-emerald-300'
                                            }`
                                            : darkMode
                                                ? 'border-gray-600 hover:border-gray-500'
                                                : 'border-gray-200 hover:border-gray-300'
                                    }`}
                                >
                                    <div className={`h-12 rounded-lg mb-2 ${themeColors[name]}`} />
                                    <p className={`font-medium capitalize ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}>{name}</p>
                                </button>
                            );
                        })}
                    </div>
                </Card>

                {/* Language */}
                <Card>
                    <div className="flex items-center gap-3 mb-4">
                        <Globe className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`} />
                        <h2 className={`font-semibold ${darkMode ? 'text-white' : 'text-gray-800'}`}>Preferred Language</h2>
                    </div>
                    <div className="relative">
                        <select
                            value={settings.preferred_language}
                            onChange={(e) => handleChange('preferred_language', e.target.value)}
                            className={`w-full px-4 py-3 pr-10 border-2 rounded-xl ${theme.focusBorder} focus:outline-none appearance-none cursor-pointer transition-colors ${
                                darkMode
                                    ? 'bg-gray-700 text-white border-gray-600 [&>option]:bg-gray-700 [&>option:checked]:bg-violet-600 [&>option:hover]:bg-violet-500'
                                    : 'bg-white text-gray-800 border-gray-200 [&>option]:bg-white [&>option:checked]:bg-violet-100 [&>option:hover]:bg-gray-100'
                            }`}
                        >
                            {LANGUAGES.map((lang) => (
                                <option key={lang.code} value={lang.code}>
                                    {lang.name}
                                </option>
                            ))}
                        </select>
                        <ChevronDown className={`absolute right-4 top-1/2 -translate-y-1/2 w-5 h-5 pointer-events-none ${darkMode ? 'text-gray-400' : 'text-gray-500'}`} />
                    </div>
                </Card>

                {/* Notifications */}
                <Card>
                    <div className="flex items-center gap-3 mb-4">
                        <Bell className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`} />
                        <h2 className={`font-semibold ${darkMode ? 'text-white' : 'text-gray-800'}`}>Notifications</h2>
                    </div>
                    <div className="space-y-4">
                        <label className="flex items-center justify-between cursor-pointer">
                            <span className={darkMode ? 'text-gray-200' : 'text-gray-700'}>Sound Effects</span>
                            <input
                                type="checkbox"
                                checked={settings.sound_enabled}
                                onChange={(e) => handleChange('sound_enabled', e.target.checked)}
                                className="w-5 h-5 rounded accent-purple-500"
                            />
                        </label>
                        <label className="flex items-center justify-between cursor-pointer">
                            <span className={darkMode ? 'text-gray-200' : 'text-gray-700'}>Email Notifications</span>
                            <input
                                type="checkbox"
                                checked={settings.email_notifications}
                                onChange={(e) => handleChange('email_notifications', e.target.checked)}
                                className="w-5 h-5 rounded accent-purple-500"
                            />
                        </label>
                    </div>
                </Card>

                {/* Privacy */}
                <Card>
                    <div className="flex items-center gap-3 mb-4">
                        <Eye className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`} />
                        <h2 className={`font-semibold ${darkMode ? 'text-white' : 'text-gray-800'}`}>Privacy</h2>
                    </div>
                    <div className="space-y-4">
                        <label className="flex items-center justify-between cursor-pointer">
                            <div>
                                <span className={darkMode ? 'text-gray-200' : 'text-gray-700'}>Public Profile</span>
                                <p className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>Allow others to see your profile</p>
                            </div>
                            <input
                                type="checkbox"
                                checked={settings.show_profile_public}
                                onChange={(e) => handleChange('show_profile_public', e.target.checked)}
                                className="w-5 h-5 rounded accent-purple-500"
                            />
                        </label>
                        <label className="flex items-center justify-between cursor-pointer">
                            <div>
                                <span className={darkMode ? 'text-gray-200' : 'text-gray-700'}>Show Statistics</span>
                                <p className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>Display your stats on leaderboards</p>
                            </div>
                            <input
                                type="checkbox"
                                checked={settings.show_stats_public}
                                onChange={(e) => handleChange('show_stats_public', e.target.checked)}
                                className="w-5 h-5 rounded accent-purple-500"
                            />
                        </label>
                    </div>
                </Card>

                {/* Messages */}
                {error && (
                    <div className={`p-4 rounded-xl border ${
                        darkMode
                            ? 'bg-red-900/30 border-red-800 text-red-400'
                            : 'bg-red-50 border-red-200 text-red-600'
                    }`}>
                        {error}
                    </div>
                )}
                {success && (
                    <div className={`p-4 rounded-xl border ${
                        darkMode
                            ? 'bg-green-900/30 border-green-800 text-green-400'
                            : 'bg-green-50 border-green-200 text-green-600'
                    }`}>
                        {success}
                    </div>
                )}

                {/* Save Button */}
                <Button
                    onClick={handleSave}
                    loading={saving}
                    icon={Save}
                    className="w-full"
                >
                    Save Settings
                </Button>
            </div>
        </Layout>
    );
}