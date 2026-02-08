import { Trophy } from '../ui/BoxIcon';
import Card from '../ui/Card';
import type { RankData } from '../../types/api';

interface RankNavigationProps {
    ranks: RankData | null;
    activeTab: string;
    onTabChange: (tab: string) => void;
}

export default function RankNavigation({
    ranks,
    activeTab,
    onTabChange,
}: RankNavigationProps) {
    const stats = [
        { id: 'global', label: 'Global', value: ranks?.global_rank },
        { id: 'monthly', label: 'Monthly', value: ranks?.monthly_rank },
        { id: 'weekly', label: 'Weekly', value: ranks?.weekly_rank },
        { id: 'daily', label: 'Daily', value: ranks?.daily_rank },
    ];

    return (
        <Card>
            <div className="flex items-center gap-3 mb-4">
                <Trophy className="w-8 h-8 text-theme-text" />
                <div>
                    <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
                        Leaderboard
                    </h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                        Click your rank to switch views
                    </p>
                </div>
            </div>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 sm:gap-3" role="tablist">
                {stats.map((stat) => {
                    const isActive = stat.id === activeTab;
                    return (
                        <button
                            key={stat.id}
                            onClick={() => onTabChange(stat.id)}
                            className={`px-3 sm:px-4 py-3 sm:py-4 rounded-xl text-center border-2 transition-all cursor-pointer ${
                                isActive
                                    ? 'border-theme-border bg-theme-gradient text-white shadow-lg'
                                    : 'border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 hover:border-gray-300 dark:hover:border-gray-500 hover:bg-gray-100 dark:hover:bg-gray-600'
                            }`}
                            role="tab"
                            aria-selected={isActive}
                        >
                            <p
                                className={`text-xs font-medium uppercase tracking-wide ${isActive ? 'text-white/80' : 'text-gray-500 dark:text-gray-400'}`}
                            >
                                {stat.label}
                            </p>
                            <p
                                className={`text-xl sm:text-2xl font-bold mt-1 ${isActive ? 'text-white' : 'text-gray-800 dark:text-white'}`}
                            >
                                {stat.value ? `#${stat.value}` : '-'}
                            </p>
                            {isActive && (
                                <div className="w-8 h-1 bg-white/50 rounded-full mx-auto mt-2" />
                            )}
                        </button>
                    );
                })}
            </div>
        </Card>
    );
}
