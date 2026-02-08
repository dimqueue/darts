import { Bell } from '../ui/BoxIcon';
import SettingsSection from './SettingsSection';
import SettingsToggle from './SettingsToggle';

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
    return (
        <SettingsSection icon={Bell} title="Notifications">
            <div className="space-y-4">
                <SettingsToggle
                    label="Sound Effects"
                    checked={soundEnabled}
                    onChange={onSoundChange}
                />
                <SettingsToggle
                    label="Email Notifications"
                    checked={emailNotifications}
                    onChange={onEmailChange}
                />
            </div>
        </SettingsSection>
    );
}
