import { useState, useEffect } from 'react';
import { X } from 'lucide-react';
import Navbar from './Navbar';
import { useTheme } from '../../contexts/ThemeContext';
import config from '../../config/env';
import { mockWords } from '../../api/mockData';

function DemoToast() {
    const [visible, setVisible] = useState(false);

    useEffect(() => {
        if (!config.useMockApi) return;

        const showTimer = setTimeout(() => setVisible(true), 500);

        const hideTimer = setTimeout(() => setVisible(false), 10500);

        return () => {
            clearTimeout(showTimer);
            clearTimeout(hideTimer);
        };
    }, []);

    if (!config.useMockApi || !visible) return null;

    return (
        <div className="fixed bottom-4 right-4 z-50 max-w-sm bg-white rounded-xl shadow-2xl border border-amber-200 overflow-hidden animate-slide-up">
            <div className="bg-amber-500 px-4 py-2 flex items-center justify-between">
                <span className="text-white font-semibold text-sm">Demo Mode</span>
                <button
                    onClick={() => setVisible(false)}
                    className="text-white/80 hover:text-white"
                >
                    <X className="w-4 h-4" />
                </button>
            </div>
            <div className="p-4">
                <p className="text-gray-600 text-sm mb-3">
                    Try guessing these secret words:
                </p>
                <div className="flex flex-wrap gap-1.5">
                    {mockWords.en.map((word) => (
                        <span key={word} className="bg-amber-100 text-amber-800 px-2 py-0.5 rounded text-xs font-mono">
                            {word}
                        </span>
                    ))}
                </div>
            </div>
        </div>
    );
}

export default function Layout({ children }) {
    const { theme } = useTheme();

    return (
        <div className={`min-h-screen bg-gradient-to-br ${theme.bgGradient}`}>
            <Navbar />
            <main className="max-w-6xl mx-auto px-4 py-8">
                {children}
            </main>
            <DemoToast />
        </div>
    );
}