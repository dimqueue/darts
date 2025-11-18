import { useState, useEffect } from 'react';
import { LogIn, User, Send, Trophy, LogOut } from 'lucide-react';
import api from './api';

// Auth Page Component
function AuthPage({ onLogin }) {
    const [isSignUp, setIsSignUp] = useState(false);
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async () => {
        setError('');
        setLoading(true);

        try {
            if (isSignUp) {
                // Sign up first
                await api.signUp(username, password);
                // Then automatically sign in
                const response = await api.signIn(username, password);
                if (response.token) {
                    onLogin(response.token, username);
                }
            } else {
                const response = await api.signIn(username, password);
                if (response.token) {
                    onLogin(response.token, username);
                }
            }
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    const handleKeyPress = (e) => {
        if (e.key === 'Enter' && username && password) {
            handleSubmit();
        }
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-green-50 to-blue-50 flex items-center justify-center p-4">
            <div className="bg-white rounded-3xl shadow-xl p-8 w-full max-w-md border-4 border-green-500">
                <div className="flex items-center justify-center mb-6">
                    <Trophy className="w-12 h-12 text-blue-400 mr-3" />
                    <h1 className="text-3xl font-bold text-gray-800">Darts Game</h1>
                </div>


                <h2 className="text-xl font-semibold text-center mb-6 text-gray-700">
                    {isSignUp ? 'Create Account' : 'Welcome Back'}
                </h2>

                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Username
                        </label>
                        <input
                            type="text"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            onKeyPress={handleKeyPress}
                            className="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-green-500 focus:outline-none"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Password
                        </label>
                        <input
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            onKeyPress={handleKeyPress}
                            className="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-green-500 focus:outline-none"
                        />
                    </div>

                    {error && (
                        <div className="text-red-500 text-sm text-center bg-red-50 p-2 rounded">
                            {error}
                        </div>
                    )}

                    <button
                        onClick={handleSubmit}
                        disabled={loading || !username || !password}
                        className="w-full bg-green-600 text-white py-3 rounded-lg font-semibold hover:bg-green-700 transition disabled:bg-gray-400 flex items-center justify-center"
                    >
                        {loading ? (
                            'Loading...'
                        ) : (
                            <>
                                <LogIn className="w-5 h-5 mr-2" />
                                {isSignUp ? 'Sign Up' : 'Sign In'}
                            </>
                        )}
                    </button>
                </div>

                <div className="mt-6 text-center">
                    <button
                        onClick={() => setIsSignUp(!isSignUp)}
                        className="text-green-600 hover:underline text-sm"
                    >
                        {isSignUp
                            ? 'Already have an account? Sign In'
                            : "Don't have an account? Sign Up"}
                    </button>
                </div>
            </div>
        </div>
    );
}

// Game Page Component
function GamePage({ username, onLogout }) {
    const [gameId, setGameId] = useState(null);
    const [guesses, setGuesses] = useState([]);
    const [inputWord, setInputWord] = useState('');
    const [loading, setLoading] = useState(false);
    const [initializing, setInitializing] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        startNewGame();
    }, []);

    const startNewGame = async () => {
        setInitializing(true);
        try {
            let language = "en";
            const game = await api.createGame(language);
            setGameId(game.id);
            setGuesses([]);
            setError('');
        } catch (err) {
            setError('Failed to start game: ' + err.message);
        } finally {
            setInitializing(false);
        }
    };

    const validateGuess = (word) => {
        const trimmed = word.trim();

        if (!trimmed) {
            return { valid: false, error: 'Please enter a word' };
        }

        if (trimmed.length < 2) {
            return { valid: false, error: 'Word must be at least 2 characters long' };
        }

        if (trimmed.length > 50) {
            return { valid: false, error: 'Word must be less than 50 characters' };
        }

        // Check for only letters (including non-English characters)
        if (!/^[\p{L}\p{M}]+$/u.test(trimmed)) {
            return { valid: false, error: 'Word must contain only letters' };
        }

        return { valid: true, word: trimmed };
    };

    const handleSubmitGuess = async () => {
        if (!gameId) return;

        const validation = validateGuess(inputWord);
        if (!validation.valid) {
            setError(validation.error);
            return;
        }

        setLoading(true);
        setError('');

        try {
            await api.createGuess(gameId, validation.word);
            // Fetch all guesses to get the updated list
            const allGuesses = await api.getAllGuessByGame(gameId);
            setGuesses(allGuesses);
            setInputWord('');
        } catch (err) {
            setError('Failed to submit guess: ' + err.message);
        } finally {
            setLoading(false);
        }
    };

    const handleKeyPress = (e) => {
        if (e.key === 'Enter' && gameId && !loading) {
            handleSubmitGuess();
        }
    };

    if (initializing) {
        return (
            <div className="min-h-screen bg-gradient-to-br from-green-50 to-blue-50 flex items-center justify-center">
                <div className="text-center">
                    <Trophy className="w-16 h-16 text-blue-400 mx-auto mb-4 animate-pulse" />
                    <p className="text-xl font-semibold text-gray-700">Starting game...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gradient-to-br from-green-50 to-blue-50 p-4">
            {/* Header */}
            <div className="max-w-2xl mx-auto mb-6 flex justify-between items-center">
                <h1 className="text-2xl font-bold text-gray-800">Darts Game</h1>
                <div className="flex items-center gap-4">
                    <div className="flex items-center gap-2 bg-white px-4 py-2 rounded-full border-2 border-blue-500">
                        <User className="w-5 h-5 text-blue-600" />
                        <span className="font-medium text-gray-700">{username}</span>
                    </div>
                    <button
                        onClick={onLogout}
                        className="bg-red-500 text-white p-2 rounded-full hover:bg-red-600 transition"
                        title="Logout"
                    >
                        <LogOut className="w-5 h-5" />
                    </button>
                </div>
            </div>

            {/* Main Game Container */}
            <div className="max-w-2xl mx-auto bg-white rounded-3xl shadow-xl p-8 border-4 border-green-500">
                {/* Input Section */}
                <div className="mb-6">
                    <label className="block text-gray-700 font-medium mb-2">
                        Input guess...
                    </label>
                    <div className="flex gap-3">
                        <input
                            type="text"
                            value={inputWord}
                            onChange={(e) => setInputWord(e.target.value)}
                            onKeyPress={handleKeyPress}
                            className="flex-1 px-4 py-3 border-2 border-gray-300 rounded-xl focus:border-green-500 focus:outline-none text-lg"
                            placeholder="Enter your word"
                            disabled={loading || !gameId}
                        />
                        <button
                            onClick={handleSubmitGuess}
                            disabled={loading || !inputWord.trim() || !gameId}
                            className="bg-green-600 text-white px-6 py-3 rounded-xl font-semibold hover:bg-green-700 transition disabled:bg-gray-400 flex items-center gap-2"
                        >
                            <Send className="w-5 h-5" />
                            Send
                        </button>
                    </div>
                </div>

                {error && (
                    <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded-lg">
                        {error}
                    </div>
                )}

                {/* Guesses Display */}
                <div className="bg-blue-50 rounded-2xl border-2 border-blue-500 p-6 min-h-[300px]">
                    {guesses.length === 0 ? (
                        <p className="text-gray-500 text-center">
                            Make your first guess!
                        </p>
                    ) : (
                        <div className="space-y-3">
                            {guesses.map((guess, index) => (
                                <div
                                    key={index}
                                    className="bg-white p-3 rounded-lg border border-gray-200"
                                >
                  <span className="text-red-600 font-bold text-lg">
                    #{index + 1}
                  </span>
                                    <span className="text-gray-700 text-lg ml-3">
                    {guess.guess_word}
                  </span>
                                    <span className="text-green-600 font-semibold ml-3">
                    Distance: {guess.distance}
                  </span>
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                {/* New Game Button */}
                <button
                    onClick={startNewGame}
                    className="mt-6 w-full bg-blue-600 text-white py-3 rounded-xl font-semibold hover:bg-blue-700 transition"
                >
                    New Game
                </button>
            </div>
        </div>
    );
}

// Main App Component
export default function App() {
    const [token, setToken] = useState(null);
    const [username, setUsername] = useState(null);

    useEffect(() => {
        // Check for existing session on mount - hydrate from localStorage
        const savedToken = localStorage.getItem('token');
        const savedUsername = localStorage.getItem('username');
        if (savedToken && savedUsername) {
            // eslint-disable-next-line react-hooks/set-state-in-effect
            setToken(savedToken);
            setUsername(savedUsername);
        }
    }, []);

    const handleLogin = (newToken, newUsername) => {
        localStorage.setItem('token', newToken);
        localStorage.setItem('username', newUsername);
        setToken(newToken);
        setUsername(newUsername);
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        setToken(null);
        setUsername(null);
    };

    return token ? (
        <GamePage username={username} onLogout={handleLogout} />
    ) : (
        <AuthPage onLogin={handleLogin} />
    );
}