import { useState, useEffect, useRef } from 'react';
import { Send, Trophy, RotateCcw, Globe } from 'lucide-react';
import api from '@/api';
import { useTheme } from '../contexts/ThemeContext';
import { useGame } from '../contexts/GameContext';
import Layout from '../components/layout/Layout';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';
import { GamePageSkeleton } from '../components/ui/Skeleton';
import {
    LANGUAGES,
    TIMEOUTS,
    GAME,
    PATTERNS,
    ERROR_CODES,
    getErrorMessage,
} from '../config/constants';

export default function GamePage() {
    const [gameId, setGameId] = useState(null);
    const [guesses, setGuesses] = useState([]);
    const [inputWord, setInputWord] = useState('');
    const [loading, setLoading] = useState(false);
    const [initializing, setInitializing] = useState(true);
    const [error, setError] = useState('');
    const [language, setLanguage] = useState('en');
    const [gameStatus, setGameStatus] = useState('in_progress');
    const initializedRef = useRef(false);
    const inputRef = useRef(null);

    const { theme, darkMode } = useTheme();
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

    const startNewGame = async (lang = language, abandonCurrent = true) => {
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
            setError('Failed to start game: ' + err.message);
        } finally {
            setInitializing(false);
            setTimeout(() => inputRef.current?.focus(), TIMEOUTS.FOCUS_DELAY);
        }
    };

    const validateGuess = (word) => {
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
            const distance = result.distance;

            if (distance === -1) {
                setError('Word not found in vocabulary. Check spelling.');
                return;
            }

            if (distance === 0) {
                setError('Word is too far from the target. Try something closer!');
                return;
            }

            if (distance === 1) {
                setGameStatus('won');
            }

            const response = await api.getAllGuessByGame(gameId);
            setGuesses(response.data || []);
            setInputWord('');
        } catch (err) {
            if (err.code === ERROR_CODES.WORD_ALREADY_USED) {
                setError(getErrorMessage(err));
            } else if (err.code === ERROR_CODES.GAME_NOT_ACTIVE) {
                setError(getErrorMessage(err));
                setGameStatus('lost');
            } else {
                setError(getErrorMessage(err));
            }
        } finally {
            setLoading(false);
            setTimeout(() => inputRef.current?.focus(), TIMEOUTS.FOCUS_DELAY_SHORT);
        }
    };

    const handleKeyPress = (e) => {
        if (e.key === 'Enter' && gameId && !loading) {
            handleSubmitGuess();
        }
    };

    const getDistanceColor = (distance) => {
        if (distance === 1) return 'bg-green-500 text-white';
        if (distance < 100) {
            return darkMode
                ? 'bg-emerald-500/30 text-emerald-200 border border-emerald-500/50'
                : 'bg-green-100 text-green-800';
        }
        if (distance < 500) {
            return darkMode
                ? 'bg-amber-500/30 text-amber-200 border border-amber-500/50'
                : 'bg-yellow-100 text-yellow-800';
        }
        if (distance < 1000) {
            return darkMode
                ? 'bg-orange-500/30 text-orange-200 border border-orange-500/50'
                : 'bg-orange-100 text-orange-800';
        }
        return darkMode
            ? 'bg-red-500/30 text-red-200 border border-red-500/50'
            : 'bg-red-100 text-red-800';
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
                {/* Game Header */}
                <Card>
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <Trophy
                                className={`w-8 h-8 ${darkMode ? theme.textColorDark : theme.textColor}`}
                            />
                            <div>
                                <h1
                                    className={`text-2xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                                >
                                    Daily Challenge
                                </h1>
                                <p
                                    className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}
                                >
                                    Guess the secret word
                                </p>
                            </div>
                        </div>
                        <div className="flex items-center gap-2">
                            <Globe
                                className={`w-5 h-5 ${darkMode ? 'text-gray-500' : 'text-gray-400'}`}
                            />
                            <select
                                value={language}
                                onChange={(e) => startNewGame(e.target.value)}
                                className={`px-3 py-2 border-2 rounded-xl ${theme.focusBorder} focus:outline-none ${
                                    darkMode
                                        ? 'bg-gray-700 text-white border-gray-600'
                                        : 'bg-white text-gray-800 border-gray-200'
                                }`}
                            >
                                {LANGUAGES.map((lang) => (
                                    <option key={lang.code} value={lang.code}>
                                        {lang.name}
                                    </option>
                                ))}
                            </select>
                        </div>
                    </div>
                </Card>

                {/* Win Message */}
                {gameStatus === 'won' && (
                    <Card className={`${theme.gradient} text-white text-center`}>
                        <Trophy className="w-12 h-12 mx-auto mb-2" />
                        <h2 className="text-2xl font-bold mb-1">Congratulations!</h2>
                        <p className="text-lg opacity-90">
                            You found the word in {guesses.filter((g) => g.distance > 0).length}{' '}
                            guesses!
                        </p>
                    </Card>
                )}

                {/* Input Section */}
                <Card>
                    <div className="flex gap-3">
                        <Input
                            ref={inputRef}
                            value={inputWord}
                            onChange={(e) => setInputWord(e.target.value)}
                            onKeyPress={handleKeyPress}
                            placeholder="Enter your guess..."
                            disabled={loading || gameStatus !== 'in_progress'}
                            className="flex-1"
                        />
                        <Button
                            onClick={handleSubmitGuess}
                            disabled={!inputWord.trim() || gameStatus !== 'in_progress'}
                            loading={loading}
                            icon={Send}
                        >
                            Send
                        </Button>
                    </div>

                    {error && (
                        <div
                            className={`mt-3 p-3 rounded-xl text-sm border ${
                                darkMode
                                    ? 'bg-red-900/30 border-red-800 text-red-400'
                                    : 'bg-red-50 border-red-200 text-red-600'
                            }`}
                        >
                            {error}
                        </div>
                    )}
                </Card>

                {/* Guesses List */}
                <Card>
                    <div className="flex items-center justify-between mb-4">
                        <h3
                            className={`font-semibold ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}
                        >
                            Your Guesses ({guesses.filter((g) => g.distance > 0).length})
                        </h3>
                        <Button
                            onClick={() => startNewGame(language)}
                            variant="outline"
                            icon={RotateCcw}
                            className="text-sm px-3 py-1"
                        >
                            New Game
                        </Button>
                    </div>

                    {guesses.filter((g) => g.distance > 0).length === 0 ? (
                        <div className="text-center py-12 text-gray-400">
                            <p>Make your first guess!</p>
                            <p className="text-sm mt-1">The closer to 0, the closer you are</p>
                        </div>
                    ) : (
                        <div className="space-y-2">
                            {guesses
                                .filter((g) => g.distance > 0)
                                .slice()
                                .sort((a, b) => a.distance - b.distance)
                                .map((guess, index) => (
                                    <div
                                        key={guess.id || index}
                                        className={`flex items-center justify-between p-3 rounded-xl ${getDistanceColor(guess.distance)}`}
                                    >
                                        <span className="font-medium">{guess.guess_word}</span>
                                        <span className="font-bold">
                                            {guess.distance === 1 ? 'FOUND!' : guess.distance}
                                        </span>
                                    </div>
                                ))}
                        </div>
                    )}
                </Card>
            </div>
        </Layout>
    );
}
