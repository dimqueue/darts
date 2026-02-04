import type { BoxIconType } from './BoxIcon';
import { useTheme } from '../../contexts/ThemeContext';

interface StatCardProps {
    label: string;
    value: string | number;
    icon?: BoxIconType;
}

export default function StatCard({ label, value, icon: Icon }: StatCardProps) {
    const { theme } = useTheme();

    return (
        <div className="rounded-xl p-4 shadow-card border bg-white dark:bg-gray-800 border-gray-100 dark:border-gray-700">
            <div className="flex items-center gap-3">
                {Icon && (
                    <div className={`p-2 rounded-lg ${theme.gradient}`}>
                        <Icon className="w-5 h-5 text-white" />
                    </div>
                )}
                <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                        {label}
                    </p>
                    <p className="text-xl font-bold text-gray-800 dark:text-white">
                        {value}
                    </p>
                </div>
            </div>
        </div>
    );
}
