import { Eye } from '../ui/BoxIcon';
import Card from '../ui/Card';
import { useTheme } from '../../contexts/ThemeContext';

interface PrivacySettingsProps {
    showProfilePublic: boolean;
    showStatsPublic: boolean;
    onProfilePublicChange: (enabled: boolean) => void;
    onStatsPublicChange: (enabled: boolean) => void;
}

export default function PrivacySettings({
    showProfilePublic,
    showStatsPublic,
    onProfilePublicChange,
    onStatsPublicChange,
}: PrivacySettingsProps) {
    const { theme, darkMode } = useTheme();

    return (
        <Card>
            <div className="flex items-center gap-3 mb-4">
                <Eye
                    className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`}
                    aria-hidden="true"
                />
                <h2 className="font-semibold text-gray-800 dark:text-white">
                    Privacy
                </h2>
            </div>
            <div className="space-y-4">
                <label className="flex items-center justify-between cursor-pointer">
                    <div>
                        <span className="text-gray-700 dark:text-gray-200">
                            Public Profile
                        </span>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                            Allow others to see your profile
                        </p>
                    </div>
                    <input
                        type="checkbox"
                        checked={showProfilePublic}
                        onChange={(e) => onProfilePublicChange(e.target.checked)}
                        className="w-5 h-5 rounded accent-purple-500"
                    />
                </label>
                <label className="flex items-center justify-between cursor-pointer">
                    <div>
                        <span className="text-gray-700 dark:text-gray-200">
                            Show Statistics
                        </span>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                            Display your stats on leaderboards
                        </p>
                    </div>
                    <input
                        type="checkbox"
                        checked={showStatsPublic}
                        onChange={(e) => onStatsPublicChange(e.target.checked)}
                        className="w-5 h-5 rounded accent-purple-500"
                    />
                </label>
            </div>
        </Card>
    );
}
