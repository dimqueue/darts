import { useState, useEffect } from 'react';
import { Settings, Palette, Globe, Bell, Eye, Save } from 'lucide-react';
import api from '../api';
import { useTheme, THEMES } from '../contexts/ThemeContext';
import Layout from '../components/layout/Layout';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';

const LANGUAGES = [
    { code: 'en', name: 'English' },
    { code: 'ua', name: 'Ukrainian' },
];

export default function SettingsPage() {
    const { themeName, setTheme, theme } = useTheme();
    const [settings, setSettings] = useState({
        preferred_language: 'en',
        theme: 'purple',
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
        loadSettings();
    }, []);

    const loadSettings = async () => {
        setLoading(true);
        try {
            const data = await api.getMySettings();
            setSettings(data);
            if (data.theme && THEMES[data.theme]) {
                setTheme(data.theme);
            }
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
    };

    const themeColors = {
        purple: 'bg-gradient-purple',
        blue: 'bg-gradient-blue',
        green: 'bg-gradient-green',
    };

    if (loading) {
        return (
            <Layout>
                <div className="flex items-center justify-center min-h-[60vh]">
                    <div className="animate-spin w-8 h-8 border-4 border-purple-500 border-t-transparent rounded-full" />
                </div>
            </Layout>
        );
    }

    return (
        <Layout>
            <div className="max-w-2xl mx-auto space-y-6">
                {/* Header */}
                <Card>
                    <div className="flex items-center gap-3">
                        <Settings className={`w-8 h-8 ${theme.textColor}`} />
                        <div>
                            <h1 className="text-2xl font-bold text-gray-800">Settings</h1>
                            <p className="text-sm text-gray-500">Customize your experience</p>
                        </div>
                    </div>
                </Card>

                {/* Theme Selection */}
                <Card>
                    <div className="flex items-center gap-3 mb-4">
                        <Palette className={`w-5 h-5 ${theme.textColor}`} />
                        <h2 className="font-semibold text-gray-800">Theme</h2>
                    </div>
                    <div className="grid grid-cols-3 gap-4">
                        {Object.keys(THEMES).map((name) => (
                            <button
                                key={name}
                                onClick={() => handleChange('theme', name)}
                                className={`p-4 rounded-xl border-2 transition-all ${
                                    settings.theme === name
                                        ? `${theme.borderColor} ring-2 ring-offset-2 ring-purple-200`
                                        : 'border-gray-200 hover:border-gray-300'
                                }`}
                            >
                                <div className={`h-12 rounded-lg mb-2 ${themeColors[name]}`} />
                                <p className="font-medium text-gray-700 capitalize">{name}</p>
                            </button>
                        ))}
                    </div>
                </Card>

                {/* Language */}
                <Card>
                    <div className="flex items-center gap-3 mb-4">
                        <Globe className={`w-5 h-5 ${theme.textColor}`} />
                        <h2 className="font-semibold text-gray-800">Preferred Language</h2>
                    </div>
                    <select
                        value={settings.preferred_language}
                        onChange={(e) => handleChange('preferred_language', e.target.value)}
                        className={`w-full px-4 py-3 border-2 rounded-xl ${theme.focusBorder} focus:outline-none`}
                    >
                        {LANGUAGES.map((lang) => (
                            <option key={lang.code} value={lang.code}>
                                {lang.name}
                            </option>
                        ))}
                    </select>
                </Card>

                {/* Notifications */}
                <Card>
                    <div className="flex items-center gap-3 mb-4">
                        <Bell className={`w-5 h-5 ${theme.textColor}`} />
                        <h2 className="font-semibold text-gray-800">Notifications</h2>
                    </div>
                    <div className="space-y-4">
                        <label className="flex items-center justify-between cursor-pointer">
                            <span className="text-gray-700">Sound Effects</span>
                            <input
                                type="checkbox"
                                checked={settings.sound_enabled}
                                onChange={(e) => handleChange('sound_enabled', e.target.checked)}
                                className="w-5 h-5 rounded accent-purple-500"
                            />
                        </label>
                        <label className="flex items-center justify-between cursor-pointer">
                            <span className="text-gray-700">Email Notifications</span>
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
                        <Eye className={`w-5 h-5 ${theme.textColor}`} />
                        <h2 className="font-semibold text-gray-800">Privacy</h2>
                    </div>
                    <div className="space-y-4">
                        <label className="flex items-center justify-between cursor-pointer">
                            <div>
                                <span className="text-gray-700">Public Profile</span>
                                <p className="text-sm text-gray-500">Allow others to see your profile</p>
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
                                <span className="text-gray-700">Show Statistics</span>
                                <p className="text-sm text-gray-500">Display your stats on leaderboards</p>
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
                    <div className="p-4 bg-red-50 border border-red-200 text-red-600 rounded-xl">
                        {error}
                    </div>
                )}
                {success && (
                    <div className="p-4 bg-green-50 border border-green-200 text-green-600 rounded-xl">
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