import { GAME_MODES } from '../../config/gameModes';
import GameModeCard from './GameModeCard';
import type { GameModeId, GameState } from '../../types/game';

const SECONDARY_MODES: GameModeId[] = ['endless', 'time-attack', 'practice'];

interface GameModesGridProps {
    getGameState: (mode: GameModeId) => GameState | null;
}

export default function GameModesGrid({ getGameState }: GameModesGridProps) {
    return (
        <section>
            <h2 className="text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500 mb-3">
                More Modes
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                {SECONDARY_MODES.map((id) => (
                    <GameModeCard
                        key={id}
                        mode={GAME_MODES[id]}
                        gameState={getGameState(id)}
                    />
                ))}
            </div>
        </section>
    );
}
