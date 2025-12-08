import { createContext, useContext, useState, useEffect } from 'react';

const GameContext = createContext(null);

const STORAGE_KEY = 'darts_game_state';

export function GameProvider({ children }) {
    const [gameState, setGameState] = useState(() => {
        const saved = localStorage.getItem(STORAGE_KEY);
        if (saved) {
            try {
                return JSON.parse(saved);
            } catch {
                localStorage.removeItem(STORAGE_KEY);
            }
        }
        return null;
    });

    useEffect(() => {
        if (gameState) {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(gameState));
        } else {
            localStorage.removeItem(STORAGE_KEY);
        }
    }, [gameState]);

    const saveGame = (state) => {
        setGameState({
            ...state,
            savedAt: new Date().toISOString(),
        });
    };

    const clearGame = () => {
        setGameState(null);
    };

    const hasActiveGame = () => {
        return gameState !== null && gameState.status === 'in_progress';
    };

    return (
        <GameContext.Provider value={{
            gameState,
            saveGame,
            clearGame,
            hasActiveGame,
        }}>
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