import { useState, useEffect } from 'react';
import { Send, Trophy, RotateCcw, Globe } from 'lucide-react';
import api from '../api';
import { useTheme } from '../contexts/ThemeContext';
import Layout from '../components/layout/Layout';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';

const LANGUAGES = [
    { code: 'en', name: 'English' },
    { code: 'ua', name: 'Ukrainian' },
];

export default function GamePage() {
    const [gameId, setGameId] = useState(null);
    const [guesses, setGuesses] = useState([]);
    const [inputWord, setInputWord] = useState('');
    const [loading, setLoading] = useState(false);
    const [initializing, setInitializing] = useState(true);
    const [error, setError] = useState('');
    const [language, setLanguage] = useState('en');
    const [gameStatus, setGameStatus] = useState('in_progress');

    const { theme } = useTheme();

    useEffect(() => {
        startNewGame();
    }, []);

    const startNewGame = async (lang = language) => {
        setInitializing(true);
        setError('');
        try {
            const game = await api.createGame(lang);
            setGameId(game.id);
            setGuesses([]);
            setGameStatus('in_progress');
            setLanguage(lang);
        } catch (err) {
            setError('Failed to start game: ' + err.message);
        } finally {
            setInitializing(false);
        }
    };

    const validateGuess = (word) => {
        const trimmed = word.trim().toLowerCase();

        if (!trimmed) {
            return { valid: false, error: 'Please enter a word' };
        }
        if (trimmed.length < 2) {
            return { valid: false, error: 'Word must be at least 2 characters' };
        }
        if (trimmed.length > 50) {
            return { valid: false, error: 'Word must be less than 50 characters' };
        }
        if (!/^[\p{L}\p{M}]+$/u.test(trimmed)) {
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

        setLoading(true);
        setError('');

        try {
            const result = await api.createGuess(gameId, validation.word);

            // Check if won (distance === 0)
            if (result.distance === 0) {
                setGameStatus('won');
            }

            const response = await api.getAllGuessByGame(gameId);
            setGuesses(response.data || []);
            setInputWord('');
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    const handleKeyPress = (e) => {
        if (e.key === 'Enter' && gameId && !loading) {
            handleSubmitGuess();
        }
    };

    const getDistanceColor = (distance) => {
        if (distance === 0) return 'bg-green-500 text-white';
        if (distance < 100) return 'bg-green-100 text-green-800';
        if (distance < 500) return 'bg-yellow-100 text-yellow-800';
        if (distance < 1000) return 'bg-orange-100 text-orange-800';
        return 'bg-red-100 text-red-800';
    };

    if (initializing) {
        return (
            <Layout>
                <div className="flex flex-col items-center justify-center min-h-[60vh]">
                    <Trophy className={`w-16 h-16 ${theme.textColor} animate-pulse mb-4`} />
                    <p className="text-xl font-semibold text-gray-700">Starting game...</p>
                </div>
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
                            <Trophy className={`w-8 h-8 ${theme.textColor}`} />
                            <div>
                                <h1 className="text-2xl font-bold text-gray-800">Daily Challenge</h1>
                                <p className="text-sm text-gray-500">Guess the secret word</p>
                            </div>
                        </div>
                        <div className="flex items-center gap-2">
                            <Globe className="w-5 h-5 text-gray-400" />
                            <select
                                value={language}
                                onChange={(e) => startNewGame(e.target.value)}
                                className={`px-3 py-2 border-2 rounded-xl ${theme.focusBorder} focus:outline-none`}
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
                        <p className="text-lg opacity-90">You found the word in {guesses.length} guesses!</p>
                    </Card>
                )}

                {/* Input Section */}
                <Card>
                    <div className="flex gap-3">
                        <Input
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
                        <div className="mt-3 p-3 bg-red-50 border border-red-200 text-red-600 rounded-xl text-sm">
                            {error}
                        </div>
                    )}
                </Card>

                {/* Guesses List */}
                <Card>
                    <div className="flex items-center justify-between mb-4">
                        <h3 className="font-semibold text-gray-700">
                            Your Guesses ({guesses.length})
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

                    {guesses.length === 0 ? (
                        <div className="text-center py-12 text-gray-400">
                            <p>Make your first guess!</p>
                            <p className="text-sm mt-1">The closer to 0, the closer you are</p>
                        </div>
                    ) : (
                        <div className="space-y-2">
                            {guesses
                                .slice()
                                .sort((a, b) => a.distance - b.distance)
                                .map((guess, index) => (
                                    <div
                                        key={guess.id || index}
                                        className={`flex items-center justify-between p-3 rounded-xl ${getDistanceColor(guess.distance)}`}
                                    >
                                        <span className="font-medium">{guess.guess_word}</span>
                                        <span className="font-bold">
                                            {guess.distance === 0 ? 'FOUND!' : guess.distance}
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