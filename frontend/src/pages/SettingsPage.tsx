import { useState, useEffect } from 'react';
import { Settings, Save } from '../components/ui/BoxIcon';
import api from '@/api';
import { useTheme } from '../contexts/ThemeContext';
import { useToast } from '../contexts/ToastContext';
import Layout from '../components/layout/Layout';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import ErrorAlert from '../components/ui/ErrorAlert';
import { SettingsSkeleton } from '../components/ui/Skeleton';
import {
    ThemeSelector,
    DarkModeToggle,
    NotificationSettings,
    PrivacySettings,
    LanguageSettings,
} from '../components/settings';
import type { SettingsData } from '../types/api';

export default function SettingsPage() {
    const { themeName, setTheme, theme, darkMode, setDarkMode } = useTheme();
    const toast = useToast();
    const [settings, setSettings] = useState<SettingsData>({
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

    useEffect(() => {
        setSettings((prev) => ({ ...prev, theme: themeName }));
    }, [themeName]);

    useEffect(() => {
        loadSettings();
    }, []);

    const loadSettings = async () => {
        setLoading(true);
        try {
            const data = await api.getMySettings();
            setSettings((prev) => ({ ...prev, ...data, theme: themeName }));
        } catch (err) {
            setError('Failed to load settings: ' + (err as Error).message);
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async () => {
        setSaving(true);
        setError('');
        try {
            await api.updateMySettings(settings);
            setTheme(settings.theme);
            toast.success('Settings saved successfully!');
        } catch (err) {
            setError('Failed to save settings: ' + (err as Error).message);
        } finally {
            setSaving(false);
        }
    };

    const handleChange = <K extends keyof SettingsData>(key: K, value: SettingsData[K]) => {
        setSettings({ ...settings, [key]: value });
        if (key === 'theme' && typeof value === 'string') {
            setTheme(value);
        }
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
            <div className="space-y-6 max-w-2xl mx-auto">
                <Card>
                    <div className="flex items-center gap-3">
                        <Settings
                            className={`w-8 h-8 ${darkMode ? theme.textColorDark : theme.textColor}`}
                            aria-hidden="true"
                        />
                        <div>
                            <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
                                Settings
                            </h1>
                            <p className="text-sm text-gray-500 dark:text-gray-400">
                                Customize your experience
                            </p>
                        </div>
                    </div>
                </Card>

                <DarkModeToggle onToggle={() => setDarkMode(!darkMode)} />

                <ThemeSelector
                    selectedTheme={settings.theme}
                    onThemeChange={(theme) => handleChange('theme', theme)}
                />

                <LanguageSettings
                    preferredLanguage={settings.preferred_language}
                    onLanguageChange={(lang) => handleChange('preferred_language', lang)}
                />

                <NotificationSettings
                    soundEnabled={settings.sound_enabled}
                    emailNotifications={settings.email_notifications}
                    onSoundChange={(enabled) => handleChange('sound_enabled', enabled)}
                    onEmailChange={(enabled) => handleChange('email_notifications', enabled)}
                />

                <PrivacySettings
                    showProfilePublic={settings.show_profile_public}
                    showStatsPublic={settings.show_stats_public}
                    onProfilePublicChange={(enabled) => handleChange('show_profile_public', enabled)}
                    onStatsPublicChange={(enabled) => handleChange('show_stats_public', enabled)}
                />

                {error && <ErrorAlert message={error} />}

                <Button onClick={handleSave} loading={saving} icon={Save} className="w-full">
                    Save Settings
                </Button>
            </div>
        </Layout>
    );
}
