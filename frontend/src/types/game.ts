import type { ComponentType } from 'react';

export type GameStatus = 'in_progress' | 'won' | 'lost' | 'abandoned';

export type GameModeId = 'daily' | 'competitive' | 'endless' | 'time-attack' | 'practice';

export type MultiplayerModeId = '1v1-duel' | 'arena';

export interface Game {
    id: number;
    language: string;
    status: GameStatus;
    started_at: string;
    finished_at?: string;
    mode?: GameModeId;
}

export interface Guess {
    id: number;
    game_id: number;
    guess_word: string;
    distance: number;
    created_at: string;
}

export interface GuessResult {
    rank: number;
    found: boolean;
    in_vocabulary: boolean;
}

export interface GameState {
    gameId: number;
    guesses: Guess[];
    language: string;
    status: GameStatus;
    mode: GameModeId;
    savedAt?: string;
}

export interface GameMode {
    id: GameModeId;
    name: string;
    description: string;
    icon: ComponentType<{ className?: string }>;
    available: boolean;
    path: string;
    gradient: string;
    rated: boolean;
    persistent: boolean;
    canContinue: boolean;
    concurrent: boolean;
    hasLeaderboard: boolean;
    comingSoon?: boolean;
}

export interface MultiplayerMode {
    id: MultiplayerModeId;
    name: string;
    description: string;
    icon: ComponentType<{ className?: string }>;
    comingSoon: boolean;
}
