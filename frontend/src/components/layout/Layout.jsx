import { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { X } from 'lucide-react';
import Navbar from './Navbar';
import { useTheme } from '../../contexts/ThemeContext';
import { useAuth } from '../../contexts/AuthContext';
import config from '../../config/env';
import { mockWords } from '../../api/mockData';

function DemoToast() {
    const [visible, setVisible] = useState(false);
    const { darkMode } = useTheme();
    const { user } = useAuth();
    const location = useLocation();

    const isGamePage = location.pathname === '/game' || location.pathname.startsWith('/game/');

    useEffect(() => {
        if (!config.useMockApi || !user || !isGamePage) {
            setVisible(false);
            return;
        }

        const timer = setTimeout(() => setVisible(true), 500);
        return () => clearTimeout(timer);
    }, [user, isGamePage]);

    const handleClose = () => {
        setVisible(false);
    };

    if (!config.useMockApi || !visible || !user || !isGamePage) return null;

    return (
        <div className={`fixed bottom-4 right-4 z-50 max-w-sm rounded-xl shadow-2xl overflow-hidden animate-slide-up ${
            darkMode ? 'bg-gray-800 border border-amber-700' : 'bg-white border border-amber-200'
        }`}>
            <div className="bg-amber-500 px-4 py-2 flex items-center justify-between">
                <span className="text-white font-semibold text-sm">Demo Mode</span>
                <button
                    onClick={handleClose}
                    className="text-white/80 hover:text-white"
                >
                    <X className="w-4 h-4" />
                </button>
            </div>
            <div className="p-4">
                <p className={`text-sm mb-3 ${darkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                    Try guessing these secret words:
                </p>

                {/* English words */}
                <p className={`text-xs font-semibold mb-1.5 ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                    English:
                </p>
                <div className="flex flex-wrap gap-1.5 mb-3">
                    {mockWords.en.map((word) => (
                        <span key={word} className={`px-2 py-0.5 rounded text-xs font-mono ${
                            darkMode ? 'bg-amber-900/50 text-amber-300' : 'bg-amber-100 text-amber-800'
                        }`}>
                            {word}
                        </span>
                    ))}
                </div>

                {/* Ukrainian words */}
                <p className={`text-xs font-semibold mb-1.5 ${darkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                    Ukrainian:
                </p>
                <div className="flex flex-wrap gap-1.5 mb-3">
                    {mockWords.ua.map((word) => (
                        <span key={word} className={`px-2 py-0.5 rounded text-xs font-mono ${
                            darkMode ? 'bg-blue-900/50 text-blue-300' : 'bg-blue-100 text-blue-800'
                        }`}>
                            {word}
                        </span>
                    ))}
                </div>

                {/* Info about changing language */}
                <p className={`text-xs mt-2 pt-2 border-t ${
                    darkMode ? 'text-gray-400 border-gray-700' : 'text-gray-500 border-gray-200'
                }`}>
                    Change language in Game Modes to try Ukrainian words
                </p>
            </div>
        </div>
    );
}

export default function Layout({ children }) {
    const { theme, darkMode } = useTheme();

    return (
        <div className={`min-h-screen bg-gradient-to-br transition-colors duration-300 ${darkMode ? theme.bgGradientDark : theme.bgGradient}`}>
            <Navbar />
            <main className="max-w-6xl mx-auto px-4 py-8">
                {children}
            </main>
            <DemoToast />
        </div>
    );
}