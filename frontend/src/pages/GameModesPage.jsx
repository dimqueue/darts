import { useNavigate } from 'react-router-dom';
import { Trophy, Clock, Zap, Infinity, Play } from 'lucide-react';
import { useTheme } from '../contexts/ThemeContext';
import Layout from '../components/layout/Layout';

const GAME_MODES = [
    {
        id: 'daily',
        name: 'Daily Challenge',
        description: 'One word per day. Compete on the leaderboard!',
        icon: Trophy,
        available: true,
        path: '/game',
        gradient: 'from-amber-400 to-orange-500',
    },
    {
        id: 'endless',
        name: 'Endless Mode',
        description: 'Play without limits. New word after each win.',
        icon: Infinity,
        available: true,
        path: '/game/endless',
        gradient: 'from-purple-400 to-indigo-500',
    },
    {
        id: 'time-attack',
        name: 'Time Attack',
        description: 'How many words in 3 minutes?',
        icon: Clock,
        available: true,
        path: '/game/time-attack',
        gradient: 'from-rose-400 to-pink-500',
    },
    {
        id: 'speed-round',
        name: 'Speed Round',
        description: 'Find the word in under 60 seconds.',
        icon: Zap,
        available: false,
        path: '/game/speed-round',
        gradient: 'from-cyan-400 to-blue-500',
    },
];

export default function GameModesPage() {
    const { theme } = useTheme();
    const navigate = useNavigate();

    const handleModeSelect = (mode) => {
        if (mode.available) {
            navigate(mode.path);
        }
    };

    return (
        <Layout>
            <div className="max-w-2xl mx-auto">
                <h1 className="text-2xl font-bold text-gray-800 mb-6 text-center">Game Modes</h1>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    {GAME_MODES.map((mode) => {
                        const Icon = mode.icon;
                        return (
                            <button
                                key={mode.id}
                                onClick={() => handleModeSelect(mode)}
                                disabled={!mode.available}
                                className={`group relative p-5 bg-white rounded-2xl text-left transition-all overflow-hidden ${
                                    mode.available
                                        ? 'hover:shadow-xl hover:-translate-y-1 cursor-pointer'
                                        : 'opacity-50'
                                }`}
                            >
                                {/* Background decoration */}
                                <div className={`absolute -right-6 -top-6 w-24 h-24 rounded-full bg-gradient-to-br ${mode.gradient} opacity-10`} />

                                {/* Icon */}
                                <div className={`relative w-12 h-12 rounded-xl flex items-center justify-center mb-4 bg-gradient-to-br ${mode.gradient} text-white shadow-lg`}>
                                    <Icon className="w-6 h-6" />
                                </div>

                                {/* Content */}
                                <h3 className="font-bold text-gray-800 mb-1">{mode.name}</h3>
                                <p className="text-sm text-gray-500 leading-relaxed">{mode.description}</p>

                                {/* Footer */}
                                <div className="mt-4 flex items-center justify-between">
                                    {mode.available ? (
                                        <span className={`inline-flex items-center gap-1 text-sm font-medium bg-gradient-to-r ${mode.gradient} bg-clip-text text-transparent`}>
                                            <Play className="w-4 h-4 text-orange-500" />
                                            Play now
                                        </span>
                                    ) : (
                                        <span className="text-xs font-medium text-gray-400 bg-gray-100 px-2 py-1 rounded-full">
                                            Coming soon
                                        </span>
                                    )}
                                </div>
                            </button>
                        );
                    })}
                </div>
            </div>
        </Layout>
    );
}