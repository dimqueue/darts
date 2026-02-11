import { useNavigate } from 'react-router-dom';
import { Play, CheckCircle } from '../ui/BoxIcon';
import type { GameMode, GameState } from '../../types/game';

interface HeroModeCardProps {
    mode: GameMode;
    gameState: GameState | null;
}

export default function HeroModeCard({ mode, gameState }: HeroModeCardProps) {
    const navigate = useNavigate();
    const Icon = mode.icon;
    const isInProgress = gameState?.status === 'in_progress';
    const isCompleted = gameState?.status === 'won';
    const guessCount = gameState?.guesses?.length ?? 0;

    return (
        <button
            onClick={() => navigate(mode.path)}
            className={`w-full text-left rounded-2xl p-5 sm:p-6 bg-gradient-to-br ${mode.gradient} text-white shadow-lg hover:shadow-xl active:scale-[0.98] hover:scale-[1.02] transition-all duration-200 group relative overflow-hidden`}
        >
            {/* Dot pattern overlay */}
            <div className="absolute inset-0 opacity-[0.07]" style={{
                backgroundImage: 'radial-gradient(circle, white 1px, transparent 1px)',
                backgroundSize: '16px 16px',
            }} />

            <div className="relative">
                <div className="flex items-start justify-between mb-4">
                    <div className="p-2.5 rounded-xl bg-white/20 backdrop-blur-sm">
                        <Icon className="w-6 h-6 text-white" />
                    </div>
                    {isInProgress && (
                        <span className="flex items-center gap-1.5 text-xs font-medium bg-white/20 backdrop-blur-sm px-2.5 py-1 rounded-full">
                            <span className="w-1.5 h-1.5 rounded-full bg-emerald-300 animate-pulse" />
                            In progress
                        </span>
                    )}
                    {isCompleted && (
                        <span className="flex items-center gap-1.5 text-xs font-medium bg-white/20 backdrop-blur-sm px-2.5 py-1 rounded-full">
                            <CheckCircle className="w-3 h-3" />
                            Done
                        </span>
                    )}
                </div>

                <h2 className="text-lg sm:text-xl font-bold mb-0.5">{mode.name}</h2>
                <p className="text-sm text-white/75 mb-4 leading-snug">{mode.description}</p>

                <div className="flex items-center gap-2 text-sm font-semibold">
                    {isCompleted ? (
                        <span className="flex items-center gap-1.5 text-white/90">
                            <CheckCircle className="w-4 h-4" />
                            Completed
                        </span>
                    ) : isInProgress ? (
                        <span className="flex items-center gap-1.5">
                            <Play className="w-4 h-4" />
                            Continue &middot; {guessCount} {guessCount === 1 ? 'guess' : 'guesses'}
                        </span>
                    ) : (
                        <span className="flex items-center gap-1.5 group-hover:translate-x-0.5 transition-transform">
                            <Play className="w-4 h-4" />
                            Play now
                        </span>
                    )}
                </div>
            </div>
        </button>
    );
}