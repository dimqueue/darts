import { createContext, useContext, useState, useEffect, useCallback, useMemo, useRef, type ReactNode } from 'react';
import { useAuth } from './AuthContext';
import { STORAGE_KEYS } from '../config/constants';
import type { GameState } from '../types/game';

function getStorageKey(userId: number | undefined): string | null {
    return userId ? `${STORAGE_KEYS.GAME_CACHE_PREFIX}${userId}` : null;
}

function getInitialGameState(userId: number | undefined): GameState | null {
    const storageKey = getStorageKey(userId);
    if (!storageKey) return null;

    const cached = localStorage.getItem(storageKey);
    if (cached) {
        try {
            return JSON.parse(cached);
        } catch {
            localStorage.removeItem(storageKey);
        }
    }
    return null;
}

interface GameContextValue {
    gameState: GameState | null;
    saveGame: (state: Omit<GameState, 'savedAt'>) => void;
    clearGame: () => void;
    hasActiveGame: () => boolean;
    getCachedGameId: () => number | null;
}

const GameContext = createContext<GameContextValue | null>(null);

interface GameProviderProps {
    children: ReactNode;
}

export function GameProvider({ children }: GameProviderProps) {
    const { user } = useAuth();
    const userId = user?.id;
    const prevUserIdRef = useRef(userId);

    const [gameState, setGameState] = useState<GameState | null>(() => getInitialGameState(userId));

    useEffect(() => {
        if (prevUserIdRef.current !== userId) {
            prevUserIdRef.current = userId;
            setGameState(getInitialGameState(userId));
        }
    }, [userId]);

    const saveGame = useCallback(
        (state: Omit<GameState, 'savedAt'>) => {
            const storageKey = getStorageKey(userId);
            const newState: GameState = {
                ...state,
                savedAt: new Date().toISOString(),
            };
            setGameState(newState);

            if (storageKey) {
                localStorage.setItem(storageKey, JSON.stringify(newState));
            }
        },
        [userId]
    );

    const clearGame = useCallback(() => {
        setGameState(null);
        const storageKey = getStorageKey(userId);
        if (storageKey) {
            localStorage.removeItem(storageKey);
        }
    }, [userId]);

    const hasActiveGame = useCallback(() => {
        return gameState !== null && gameState.status === 'in_progress';
    }, [gameState]);

    const getCachedGameId = useCallback(() => {
        return gameState?.gameId || null;
    }, [gameState]);

    const value = useMemo<GameContextValue>(
        () => ({ gameState, saveGame, clearGame, hasActiveGame, getCachedGameId }),
        [gameState, saveGame, clearGame, hasActiveGame, getCachedGameId]
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
