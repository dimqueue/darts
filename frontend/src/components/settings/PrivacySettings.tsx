import { Eye } from '../ui/BoxIcon';
import SettingsSection from './SettingsSection';
import SettingsToggle from './SettingsToggle';

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
    return (
        <SettingsSection icon={Eye} title="Privacy">
            <div className="space-y-4">
                <SettingsToggle
                    label="Public Profile"
                    description="Allow others to see your profile"
                    checked={showProfilePublic}
                    onChange={onProfilePublicChange}
                />
                <SettingsToggle
                    label="Show Statistics"
                    description="Display your stats on leaderboards"
                    checked={showStatsPublic}
                    onChange={onStatsPublicChange}
                />
            </div>
        </SettingsSection>
    );
}
