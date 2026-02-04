import { useState, useEffect, useRef } from 'react';
import api from '@/api';
import { useGame } from '../contexts/GameContext';
import Layout from '../components/layout/Layout';
import { GamePageSkeleton } from '../components/ui/Skeleton';
import { GameHeader, GuessInput, GuessList, WinMessage } from '../components/game';
import { TIMEOUTS, GAME, PATTERNS, ERROR_CODES, getErrorMessage } from '../config/constants';
import type { Guess, GameStatus } from '../types/game';

export default function GamePage() {
    const [gameId, setGameId] = useState<number | null>(null);
    const [guesses, setGuesses] = useState<Guess[]>([]);
    const [inputWord, setInputWord] = useState('');
    const [loading, setLoading] = useState(false);
    const [initializing, setInitializing] = useState(true);
    const [error, setError] = useState('');
    const [language, setLanguage] = useState('en');
    const [gameStatus, setGameStatus] = useState<GameStatus>('in_progress');
    const initializedRef = useRef(false);
    const inputRef = useRef<HTMLInputElement>(null);

    const { saveGame } = useGame();

    useEffect(() => {
        if (gameId && !initializing) {
            saveGame({ gameId, guesses, language, status: gameStatus });
        }
    }, [gameId, guesses, language, gameStatus, initializing, saveGame]);

    useEffect(() => {
        if (initializedRef.current) return;
        initializedRef.current = true;

        const loadActiveGame = async () => {
            try {
                const response = await api.getActiveGame();
                const activeGame = response.data;

                if (activeGame) {
                    setGameId(activeGame.id);
                    setLanguage(activeGame.language || 'en');
                    setGameStatus(activeGame.status);

                    const guessesResponse = await api.getAllGuessByGame(activeGame.id);
                    setGuesses(guessesResponse.data || []);
                    setInitializing(false);
                    setTimeout(() => inputRef.current?.focus(), TIMEOUTS.FOCUS_DELAY);
                } else {
                    await startNewGame('en', false);
                }
            } catch (err) {
                console.error('Failed to load active game:', err);
                await startNewGame('en', false);
            }
        };

        loadActiveGame();
    }, []);

    const startNewGame = async (lang: string = language, abandonCurrent: boolean = true) => {
        setInitializing(true);
        setError('');
        try {
            if (abandonCurrent && gameId && gameStatus === 'in_progress') {
                await api.abandonGame(gameId).catch(() => {});
            }

            const game = await api.createGame(lang);
            setGameId(game.id);
            setGuesses([]);
            setGameStatus('in_progress');
            setLanguage(lang);
        } catch (err) {
            setError('Failed to start game: ' + (err as Error).message);
        } finally {
            setInitializing(false);
            setTimeout(() => inputRef.current?.focus(), TIMEOUTS.FOCUS_DELAY);
        }
    };

    const validateGuess = (
        word: string
    ): { valid: false; error: string } | { valid: true; word: string } => {
        const trimmed = word.trim().toLowerCase();

        if (!trimmed) {
            return { valid: false, error: 'Please enter a word' };
        }
        if (trimmed.length < GAME.MIN_WORD_LENGTH) {
            return {
                valid: false,
                error: `Word must be at least ${GAME.MIN_WORD_LENGTH} characters`,
            };
        }
        if (trimmed.length > GAME.MAX_WORD_LENGTH) {
            return {
                valid: false,
                error: `Word must be less than ${GAME.MAX_WORD_LENGTH} characters`,
            };
        }
        if (!PATTERNS.WORD.test(trimmed)) {
            return { valid: false, error: 'Word must contain only letters' };
        }
        return { valid: true, word: trimmed };
    };

    const handleSubmitGuess = async () => {
        if (!gameId || gameStatus !== 'in_progress') return;

        const validation = validateGuess(inputWord);
        if (!validation.valid) {
            setError(validation.error);
            return;
        }

        const alreadyGuessed = guesses.some(
            (g) => g.guess_word.toLowerCase() === validation.word.toLowerCase()
        );
        if (alreadyGuessed) {
            setError('You already guessed this word');
            return;
        }

        setLoading(true);
        setError('');

        try {
            const result = await api.createGuess(gameId, validation.word);

            if (!result.in_vocabulary) {
                setError('Word not found in vocabulary. Check spelling.');
                return;
            }

            if (result.found) {
                setGameStatus('won');
            }

            const response = await api.getAllGuessByGame(gameId);
            setGuesses(response.data || []);
            setInputWord('');
        } catch (err) {
            const appError = err as { code?: string; message?: string };
            if (appError.code === ERROR_CODES.WORD_ALREADY_USED) {
                setError(getErrorMessage(appError));
            } else if (appError.code === ERROR_CODES.GAME_NOT_ACTIVE) {
                setError(getErrorMessage(appError));
                setGameStatus('lost');
            } else {
                setError(getErrorMessage(appError));
            }
        } finally {
            setLoading(false);
            setTimeout(() => inputRef.current?.focus(), TIMEOUTS.FOCUS_DELAY_SHORT);
        }
    };

    if (initializing) {
        return (
            <Layout>
                <GamePageSkeleton />
            </Layout>
        );
    }

    return (
        <Layout>
            <div className="max-w-2xl mx-auto space-y-6">
                <GameHeader language={language} onLanguageChange={startNewGame} />

                {gameStatus === 'won' && (
                    <WinMessage guessCount={guesses.filter((g) => g.distance > 0).length} />
                )}

                <GuessInput
                    ref={inputRef}
                    value={inputWord}
                    onChange={setInputWord}
                    onSubmit={handleSubmitGuess}
                    loading={loading}
                    disabled={gameStatus !== 'in_progress'}
                    error={error}
                />

                <GuessList guesses={guesses} onNewGame={() => startNewGame(language)} />
            </div>
        </Layout>
    );
}
