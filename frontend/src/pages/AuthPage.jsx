import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { LogIn, Trophy } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import api from '@/api';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';

export default function AuthPage() {
    const [isSignUp, setIsSignUp] = useState(false);
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const { login } = useAuth();
    const { theme, darkMode } = useTheme();
    const navigate = useNavigate();

    const handleSubmit = async () => {
        setError('');
        setLoading(true);

        try {
            if (isSignUp) {
                await api.signUp(username, password);
                const response = await api.signIn(username, password);
                if (response.token) {
                    login(response.token, { username, id: response.user_id });
                    navigate('/game');
                }
            } else {
                const response = await api.signIn(username, password);
                if (response.token) {
                    login(response.token, { username, id: response.user_id });
                    navigate('/game');
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
        <div
            className={`min-h-screen bg-gradient-to-br transition-colors duration-300 ${darkMode ? theme.bgGradientDark : theme.bgGradient} flex items-center justify-center p-4`}
        >
            <Card className="w-full max-w-md" padding="p-8">
                <div className="flex items-center justify-center mb-6">
                    <Trophy
                        className={`w-12 h-12 ${darkMode ? theme.textColorDark : theme.textColor} mr-3`}
                    />
                    <h1
                        className={`text-3xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}
                    >
                        Darts Game
                    </h1>
                </div>

                <h2
                    className={`text-xl font-semibold text-center mb-6 ${darkMode ? 'text-gray-300' : 'text-gray-700'}`}
                >
                    {isSignUp ? 'Create Account' : 'Welcome Back'}
                </h2>

                <div className="space-y-4">
                    <Input
                        label="Username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        onKeyPress={handleKeyPress}
                        placeholder="Enter your username"
                    />

                    <Input
                        label="Password"
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        onKeyPress={handleKeyPress}
                        placeholder="Enter your password"
                    />

                    {error && (
                        <div
                            className={`text-sm text-center p-3 rounded-xl ${
                                darkMode
                                    ? 'bg-red-900/30 text-red-400 border border-red-800'
                                    : 'bg-red-50 text-red-500'
                            }`}
                        >
                            {error}
                        </div>
                    )}

                    <Button
                        onClick={handleSubmit}
                        disabled={!username || !password}
                        loading={loading}
                        icon={LogIn}
                        className="w-full"
                    >
                        {isSignUp ? 'Sign Up' : 'Sign In'}
                    </Button>
                </div>

                <div className="mt-6 text-center">
                    <button
                        onClick={() => setIsSignUp(!isSignUp)}
                        className={`${darkMode ? theme.textColorDark : theme.textColor} hover:underline text-sm`}
                    >
                        {isSignUp
                            ? 'Already have an account? Sign In'
                            : "Don't have an account? Sign Up"}
                    </button>
                </div>
            </Card>
        </div>
    );
}
