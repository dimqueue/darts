import { useState, useEffect, useRef } from 'react';
import api from '@/api';
import { useAuth } from '../contexts/AuthContext';
import Layout from '../components/layout/Layout';
import { ProfileSkeleton } from '../components/ui/Skeleton';
import {
    ProfileHeader,
    StatsGrid,
    DetailedStats,
    LanguageStats,
} from '../components/profile';
import type { ProfileData, StatisticsData, LanguageStatsData, EditFormData } from '../types/api';

export default function ProfilePage() {
    const { user, updateUser } = useAuth();
    const [profile, setProfile] = useState<ProfileData | null>(null);
    const [stats, setStats] = useState<StatisticsData | null>(null);
    const [languageStats, setLanguageStats] = useState<LanguageStatsData[]>([]);
    const [loading, setLoading] = useState(true);
    const [editing, setEditing] = useState(false);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');
    const mountedRef = useRef(true);

    const [editForm, setEditForm] = useState<EditFormData>({
        bio: '',
        country_code: '',
    });

    useEffect(() => {
        mountedRef.current = true;
        loadProfile();
        return () => {
            mountedRef.current = false;
        };
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

            if (!mountedRef.current) return;

            setProfile(profileData);
            setStats(statsData);
            setLanguageStats(langStatsData?.data || []);
            setEditForm({
                bio: profileData.bio || '',
                country_code: profileData.country_code || '',
            });
        } catch (err) {
            if (!mountedRef.current) return;
            setError('Failed to load profile: ' + (err as Error).message);
        } finally {
            if (mountedRef.current) {
                setLoading(false);
            }
        }
    };

    const handleSave = async () => {
        setSaving(true);
        setError('');
        try {
            await api.updateMyProfile(editForm);
            if (!mountedRef.current) return;
            setProfile({ ...profile!, ...editForm });
            updateUser({ ...user!, ...editForm });
            setEditing(false);
        } catch (err) {
            if (!mountedRef.current) return;
            setError('Failed to save: ' + (err as Error).message);
        } finally {
            if (mountedRef.current) {
                setSaving(false);
            }
        }
    };

    if (loading) {
        return (
            <Layout>
                <ProfileSkeleton />
            </Layout>
        );
    }

    return (
        <Layout>
            <div className="space-y-6">
                <ProfileHeader
                    profile={profile}
                    username={user?.username}
                    editing={editing}
                    saving={saving}
                    error={error}
                    editForm={editForm}
                    onToggleEdit={() => setEditing(!editing)}
                    onSave={handleSave}
                    onEditFormChange={setEditForm}
                />

                <StatsGrid stats={stats} />

                <DetailedStats stats={stats} />

                <LanguageStats languageStats={languageStats} />
            </div>
        </Layout>
    );
}
