import Card from '../ui/Card';
import type { StatisticsData } from '../../types/api';

interface DetailedStatsProps {
    stats: StatisticsData | null;
}

export default function DetailedStats({ stats }: DetailedStatsProps) {
    const statItems = [
        { label: 'Current Streak', value: stats?.current_win_streak || 0 },
        { label: 'Avg. Guesses', value: stats?.average_guesses?.toFixed(1) || '0.0' },
        { label: 'Total Score', value: stats?.total_score || 0 },
        {
            label: 'Fastest Win',
            value: stats?.fastest_win_seconds ? `${stats.fastest_win_seconds}s` : '-',
        },
        { label: 'Fewest Guesses', value: stats?.fewest_guesses_win || '-' },
        { label: 'Total Guesses', value: stats?.total_guesses || 0 },
    ];

    return (
        <Card>
            <h2 className="text-lg font-semibold mb-4 text-gray-800 dark:text-white">
                Detailed Statistics
            </h2>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                {statItems.map((item) => (
                    <div key={item.label}>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                            {item.label}
                        </p>
                        <p className="text-xl font-bold text-gray-800 dark:text-white">
                            {item.value}
                        </p>
                    </div>
                ))}
            </div>
        </Card>
    );
}
