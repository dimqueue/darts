import { useNavigate } from 'react-router-dom';
import type { GameMode, GameState } from '../../types/game';

interface GameModeCardProps {
    mode: GameMode;
    gameState: GameState | null;
}

export default function GameModeCard({ mode, gameState }: GameModeCardProps) {
    const navigate = useNavigate();
    const Icon = mode.icon;
    const isActive = gameState?.status === 'in_progress';
    const disabled = !mode.available || mode.comingSoon;

    return (
        <button
            onClick={() => !disabled && navigate(mode.path)}
            disabled={disabled}
            className={`w-full text-left rounded-xl p-4 border transition-all duration-200 ${
                disabled
                    ? 'opacity-50 cursor-not-allowed bg-gray-50 dark:bg-gray-800/50 border-gray-100 dark:border-gray-700'
                    : 'bg-white dark:bg-gray-800 border-gray-100 dark:border-gray-700 shadow-card hover:shadow-lg active:scale-[0.98] hover:scale-[1.02]'
            }`}
        >
            <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg bg-gradient-to-br ${mode.gradient} shrink-0`}>
                    <Icon className="w-4 h-4 text-white" />
                </div>
                <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                        <span className="font-semibold text-sm text-gray-800 dark:text-white">
                            {mode.name}
                        </span>
                        {mode.comingSoon && (
                            <span className="text-[10px] px-1.5 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-400 dark:text-gray-500 font-medium uppercase tracking-wide">
                                Soon
                            </span>
                        )}
                        {isActive && (
                            <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                        )}
                    </div>
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5 truncate">
                        {mode.description}
                    </p>
                </div>
            </div>
        </button>
    );
}
