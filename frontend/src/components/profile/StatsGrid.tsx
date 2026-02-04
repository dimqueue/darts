import { Trophy, Target, TrendingUp, Flame } from '../ui/BoxIcon';
import StatCard from '../ui/StatCard';
import type { StatisticsData } from '../../types/api';

interface StatsGridProps {
    stats: StatisticsData | null;
}

export default function StatsGrid({ stats }: StatsGridProps) {
    const winRate =
        stats && stats.total_games > 0
            ? ((stats.total_wins / stats.total_games) * 100).toFixed(1)
            : '0';

    return (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <StatCard label="Total Games" value={stats?.total_games || 0} icon={Target} />
            <StatCard label="Wins" value={stats?.total_wins || 0} icon={Trophy} />
            <StatCard label="Win Rate" value={`${winRate}%`} icon={TrendingUp} />
            <StatCard
                label="Best Streak"
                value={stats?.best_win_streak || 0}
                icon={Flame}
            />
        </div>
    );
}
