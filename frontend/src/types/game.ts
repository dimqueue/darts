import type { ComponentType } from 'react';

export type GameStatus = 'in_progress' | 'won' | 'lost' | 'abandoned';

export interface Game {
    id: number;
    language: string;
    status: GameStatus;
    started_at: string;
    finished_at?: string;
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
    savedAt?: string;
}

export interface GameMode {
    id: string;
    name: string;
    description: string;
    icon: ComponentType<{ className?: string }>;
    available: boolean;
    path: string;
    gradient: string;
}
