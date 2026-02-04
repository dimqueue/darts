import { Bell } from '../ui/BoxIcon';
import Card from '../ui/Card';
import { useTheme } from '../../contexts/ThemeContext';

interface NotificationSettingsProps {
    soundEnabled: boolean;
    emailNotifications: boolean;
    onSoundChange: (enabled: boolean) => void;
    onEmailChange: (enabled: boolean) => void;
}

export default function NotificationSettings({
    soundEnabled,
    emailNotifications,
    onSoundChange,
    onEmailChange,
}: NotificationSettingsProps) {
    const { theme, darkMode } = useTheme();

    return (
        <Card>
            <div className="flex items-center gap-3 mb-4">
                <Bell
                    className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`}
                    aria-hidden="true"
                />
                <h2 className="font-semibold text-gray-800 dark:text-white">
                    Notifications
                </h2>
            </div>
            <div className="space-y-4">
                <label className="flex items-center justify-between cursor-pointer">
                    <span className="text-gray-700 dark:text-gray-200">
                        Sound Effects
                    </span>
                    <input
                        type="checkbox"
                        checked={soundEnabled}
                        onChange={(e) => onSoundChange(e.target.checked)}
                        className="w-5 h-5 rounded accent-purple-500"
                    />
                </label>
                <label className="flex items-center justify-between cursor-pointer">
                    <span className="text-gray-700 dark:text-gray-200">
                        Email Notifications
                    </span>
                    <input
                        type="checkbox"
                        checked={emailNotifications}
                        onChange={(e) => onEmailChange(e.target.checked)}
                        className="w-5 h-5 rounded accent-purple-500"
                    />
                </label>
            </div>
        </Card>
    );
}
