import { useState, useEffect } from 'react';
import api from '@/api';
import { useAuth } from '../contexts/AuthContext';
import { useGame } from '../contexts/GameContext';
import Layout from '../components/layout/Layout';
import { HomePageSkeleton } from '../components/ui/Skeleton';
import { HeroModeCard, GameModesGrid } from '../components/home';
import { GAME_MODES } from '../config/gameModes';
import { Crown, Flame } from '../components/ui/BoxIcon';
import type { StatisticsData, RankData } from '../types/api';

export default function HomePage() {
    const { user } = useAuth();
    const { getGameState } = useGame();
    const [stats, setStats] = useState<StatisticsData | null>(null);
    const [rank, setRank] = useState<RankData | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        Promise.all([
            api.getMyStatistics().catch(() => null),
            api.getMyRank().catch(() => null),
        ]).then(([statsData, rankData]) => {
            setStats(statsData);
            setRank(rankData);
            setLoading(false);
        });
    }, []);

    if (loading) {
        return (
            <Layout>
                <HomePageSkeleton />
            </Layout>
        );
    }

    const displayName = user?.name || user?.username || 'Player';

    return (
        <Layout>
            <div className="max-w-4xl mx-auto space-y-6">
                {/* Welcome banner */}
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
                            Welcome back, {displayName}
                        </h1>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                            Ready for a challenge?
                        </p>
                    </div>
                    <div className="flex items-center gap-2">
                        {rank?.global_rank != null && rank.global_rank > 0 && (
                            <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-theme-light-bg text-theme-text text-sm font-medium">
                                <Crown className="w-4 h-4" />
                                #{rank.global_rank}
                            </div>
                        )}
                        {stats?.current_win_streak != null && stats.current_win_streak > 0 && (
                            <div className="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-orange-50 dark:bg-orange-900/20 text-orange-600 dark:text-orange-400 text-sm font-medium">
                                <Flame className="w-4 h-4" />
                                {stats.current_win_streak}
                            </div>
                        )}
                    </div>
                </div>

                {/* Hero mode cards */}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <HeroModeCard
                        mode={GAME_MODES.daily}
                        gameState={getGameState('daily')}
                    />
                    <HeroModeCard
                        mode={GAME_MODES.competitive}
                        gameState={getGameState('competitive')}
                    />
                </div>

                {/* Secondary modes */}
                <GameModesGrid getGameState={getGameState} />
            </div>
        </Layout>
    );
}
