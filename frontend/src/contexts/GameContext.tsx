import { createContext, useContext, useState, useEffect, useCallback, useMemo, useRef, type ReactNode } from 'react';
import { useAuth } from './AuthContext';
import { STORAGE_KEYS } from '../config/constants';
import { ALL_MODE_IDS } from '../config/gameModes';
import type { GameModeId, GameState } from '../types/game';

function getModeStorageKey(userId: number | undefined, mode: GameModeId): string | null {
    return userId ? `${STORAGE_KEYS.GAME_CACHE_PREFIX}${userId}_${mode}` : null;
}

function getOldStorageKey(userId: number | undefined): string | null {
    return userId ? `${STORAGE_KEYS.GAME_CACHE_PREFIX}${userId}` : null;
}

function getInitialGameStates(userId: number | undefined): Map<GameModeId, GameState> {
    const states = new Map<GameModeId, GameState>();
    if (!userId) return states;

    const oldKey = getOldStorageKey(userId);
    const competitiveKey = getModeStorageKey(userId, 'competitive');
    if (oldKey && competitiveKey) {
        const oldCached = localStorage.getItem(oldKey);
        if (oldCached && !localStorage.getItem(competitiveKey)) {
            try {
                const parsed = JSON.parse(oldCached);
                parsed.mode = 'competitive';
                localStorage.setItem(competitiveKey, JSON.stringify(parsed));
            } catch { /* ignore corrupt data */ }
            localStorage.removeItem(oldKey);
        }
    }

    for (const mode of ALL_MODE_IDS) {
        const key = getModeStorageKey(userId, mode);
        if (!key) continue;
        const cached = localStorage.getItem(key);
        if (cached) {
            try {
                const parsed = JSON.parse(cached) as GameState;
                parsed.mode = mode;
                states.set(mode, parsed);
            } catch {
                localStorage.removeItem(key);
            }
        }
    }

    return states;
}

interface GameContextValue {
    getGameState: (mode: GameModeId) => GameState | null;
    saveGame: (state: Omit<GameState, 'savedAt'>) => void;
    clearGame: (mode: GameModeId) => void;
    hasActiveGame: (mode: GameModeId) => boolean;
    getCachedGameId: (mode: GameModeId) => number | null;
    getActiveGameModes: () => GameModeId[];
    clearAllGames: () => void;
}

const GameContext = createContext<GameContextValue | null>(null);

interface GameProviderProps {
    children: ReactNode;
}

export function GameProvider({ children }: GameProviderProps) {
    const { user } = useAuth();
    const userId = user?.id;
    const prevUserIdRef = useRef(userId);

    const [gameStates, setGameStates] = useState<Map<GameModeId, GameState>>(
        () => getInitialGameStates(userId)
    );

    useEffect(() => {
        if (prevUserIdRef.current !== userId) {
            prevUserIdRef.current = userId;
            setGameStates(getInitialGameStates(userId));
        }
    }, [userId]);

    const getGameState = useCallback(
        (mode: GameModeId): GameState | null => gameStates.get(mode) ?? null,
        [gameStates]
    );

    const saveGame = useCallback(
        (state: Omit<GameState, 'savedAt'>) => {
            const mode = state.mode;
            const newState: GameState = {
                ...state,
                savedAt: new Date().toISOString(),
            };
            setGameStates((prev) => {
                const next = new Map(prev);
                next.set(mode, newState);
                return next;
            });

            const key = getModeStorageKey(userId, mode);
            if (key) {
                localStorage.setItem(key, JSON.stringify(newState));
            }
        },
        [userId]
    );

    const clearGame = useCallback(
        (mode: GameModeId) => {
            setGameStates((prev) => {
                const next = new Map(prev);
                next.delete(mode);
                return next;
            });
            const key = getModeStorageKey(userId, mode);
            if (key) {
                localStorage.removeItem(key);
            }
        },
        [userId]
    );

    const hasActiveGame = useCallback(
        (mode: GameModeId): boolean => {
            const state = gameStates.get(mode);
            return state !== null && state !== undefined && state.status === 'in_progress';
        },
        [gameStates]
    );

    const getCachedGameId = useCallback(
        (mode: GameModeId): number | null => gameStates.get(mode)?.gameId ?? null,
        [gameStates]
    );

    const getActiveGameModes = useCallback(
        (): GameModeId[] =>
            ALL_MODE_IDS.filter((mode) => {
                const state = gameStates.get(mode);
                return state && state.status === 'in_progress';
            }),
        [gameStates]
    );

    const clearAllGames = useCallback(() => {
        setGameStates(new Map());
        if (!userId) return;
        for (const mode of ALL_MODE_IDS) {
            const key = getModeStorageKey(userId, mode);
            if (key) localStorage.removeItem(key);
        }
        // Also clean up old key if it somehow still exists
        const oldKey = getOldStorageKey(userId);
        if (oldKey) localStorage.removeItem(oldKey);
    }, [userId]);

    const value = useMemo<GameContextValue>(
        () => ({
            getGameState,
            saveGame,
            clearGame,
            hasActiveGame,
            getCachedGameId,
            getActiveGameModes,
            clearAllGames,
        }),
        [getGameState, saveGame, clearGame, hasActiveGame, getCachedGameId, getActiveGameModes, clearAllGames]
    );

    return (
        <GameContext.Provider value={value}>
            {children}
        </GameContext.Provider>
    );
}

export function useGame(): GameContextValue {
    const context = useContext(GameContext);
    if (!context) {
        throw new Error('useGame must be used within a GameProvider');
    }
    return context;
}
