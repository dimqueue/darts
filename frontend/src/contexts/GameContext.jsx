import { createContext, useContext, useState, useEffect, useCallback, useRef } from 'react';
import { useAuth } from './AuthContext';

const GameContext = createContext(null);

const STORAGE_KEY_PREFIX = 'darts_game_cache_';

function getStorageKey(userId) {
    return userId ? `${STORAGE_KEY_PREFIX}${userId}` : null;
}

function getInitialGameState(userId) {
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

export function GameProvider({ children }) {
    const { user } = useAuth();
    const userId = user?.id;
    const prevUserIdRef = useRef(userId);

    const [gameState, setGameState] = useState(() => getInitialGameState(userId));

    useEffect(() => {
        if (prevUserIdRef.current !== userId) {
            prevUserIdRef.current = userId;
            setGameState(getInitialGameState(userId));
        }
    }, [userId]);

    const saveGame = useCallback(
        (state) => {
            const storageKey = getStorageKey(userId);
            const newState = {
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

    return (
        <GameContext.Provider
            value={{
                gameState,
                saveGame,
                clearGame,
                hasActiveGame,
                getCachedGameId,
            }}
        >
            {children}
        </GameContext.Provider>
    );
}

export function useGame() {
    const context = useContext(GameContext);
    if (!context) {
        throw new Error('useGame must be used within a GameProvider');
    }
    return context;
}
