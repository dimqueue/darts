import type { ReactNode } from 'react';
import type { BoxIconType } from '../ui/BoxIcon';
import Card from '../ui/Card';

interface SettingsSectionProps {
    icon: BoxIconType;
    title: string;
    children: ReactNode;
}

export default function SettingsSection({ icon: Icon, title, children }: SettingsSectionProps) {
    return (
        <Card>
            <div className="flex items-center gap-3 mb-4">
                <Icon className="w-5 h-5 text-theme-text" aria-hidden="true" />
                <h2 className="font-semibold text-gray-800 dark:text-white">{title}</h2>
            </div>
            {children}
        </Card>
    );
}
