import { User, Edit2, Save } from 'lucide-react';
import Card from '../ui/Card';
import Button from '../ui/Button';
import Input from '../ui/Input';
import ErrorAlert from '../ui/ErrorAlert';
import { useTheme } from '../../contexts/ThemeContext';
import type { ProfileData } from '../../types/api';

interface EditFormData {
    bio: string;
    country_code: string;
}

interface ProfileHeaderProps {
    profile: ProfileData | null;
    username: string | undefined;
    editing: boolean;
    saving: boolean;
    error: string;
    editForm: EditFormData;
    onToggleEdit: () => void;
    onSave: () => void;
    onEditFormChange: (form: EditFormData) => void;
}

export default function ProfileHeader({
    profile,
    username,
    editing,
    saving,
    error,
    editForm,
    onToggleEdit,
    onSave,
    onEditFormChange,
}: ProfileHeaderProps) {
    const { theme, darkMode } = useTheme();

    return (
        <Card>
            <div className="flex items-start justify-between">
                <div className="flex items-center gap-4">
                    <div
                        className={`w-20 h-20 ${theme.gradient} rounded-full flex items-center justify-center`}
                    >
                        <User className="w-10 h-10 text-white" aria-hidden="true" />
                    </div>
                    <div>
                        <h1
                            className={`text-2xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                        >
                            {profile?.name || username}
                        </h1>
                        <p className={darkMode ? 'text-gray-400' : 'text-gray-500'}>
                            @{profile?.username || username}
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
                    onClick={onToggleEdit}
                    variant="outline"
                    icon={editing ? undefined : Edit2}
                >
                    {editing ? 'Cancel' : 'Edit'}
                </Button>
            </div>

            {editing ? (
                <div className="mt-6 space-y-4">
                    <Input
                        label="Bio"
                        value={editForm.bio}
                        onChange={(e) =>
                            onEditFormChange({ ...editForm, bio: e.target.value })
                        }
                        placeholder="Tell us about yourself..."
                    />
                    <Input
                        label="Country Code"
                        value={editForm.country_code}
                        onChange={(e) =>
                            onEditFormChange({
                                ...editForm,
                                country_code: e.target.value.toUpperCase().slice(0, 2),
                            })
                        }
                        placeholder="US, GB, DE..."
                    />
                    <Button onClick={onSave} loading={saving} icon={Save}>
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

            {error && <ErrorAlert message={error} className="mt-4" />}
        </Card>
    );
}
