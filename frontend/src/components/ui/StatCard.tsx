import type { LucideIcon } from 'lucide-react';
import { useTheme } from '../../contexts/ThemeContext';

interface StatCardProps {
    label: string;
    value: string | number;
    icon?: LucideIcon;
}

export default function StatCard({ label, value, icon: Icon }: StatCardProps) {
    const { theme, darkMode } = useTheme();

    return (
        <div
            className={`rounded-xl p-4 shadow-card border transition-colors duration-300 ${
                darkMode ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'
            }`}
        >
            <div className="flex items-center gap-3">
                {Icon && (
                    <div className={`p-2 rounded-lg ${theme.gradient}`}>
                        <Icon className="w-5 h-5 text-white" />
                    </div>
                )}
                <div>
                    <p className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                        {label}
                    </p>
                    <p className={`text-xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}>
                        {value}
                    </p>
                </div>
            </div>
        </div>
    );
}
