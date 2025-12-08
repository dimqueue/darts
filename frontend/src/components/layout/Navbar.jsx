import { NavLink } from 'react-router-dom';
import { Target, Crosshair, Layers, Crown, User, Settings, LogOut, Palette, Moon, Sun } from 'lucide-react';
import { useAuth } from '../../contexts/AuthContext';
import { useGame } from '../../contexts/GameContext';
import { useTheme, THEMES } from '../../contexts/ThemeContext';
import { useState, useRef, useEffect } from 'react';

export default function Navbar() {
    const { user, logout } = useAuth();
    const { clearGame } = useGame();
    const { themeName, setTheme, theme, darkMode, toggleDarkMode } = useTheme();
    const [showThemeMenu, setShowThemeMenu] = useState(false);
    const [showUserMenu, setShowUserMenu] = useState(false);
    const themeMenuRef = useRef(null);
    const userMenuRef = useRef(null);

    useEffect(() => {
        function handleClickOutside(event) {
            if (themeMenuRef.current && !themeMenuRef.current.contains(event.target)) {
                setShowThemeMenu(false);
            }
            if (userMenuRef.current && !userMenuRef.current.contains(event.target)) {
                setShowUserMenu(false);
            }
        }
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const navLinks = [
        { to: '/game', label: 'Play', icon: Crosshair },
        { to: '/modes', label: 'Modes', icon: Layers },
        { to: '/leaderboard', label: 'Leaderboard', icon: Crown },
    ];

    const themeColors = {
        purple: 'bg-purple-500',
        blue: 'bg-blue-500',
        green: 'bg-green-500',
    };

    return (
        <nav className={`sticky top-0 z-50 shadow-md transition-colors duration-300 ${darkMode ? 'bg-gray-900' : 'bg-white'}`}>
            <div className="max-w-6xl mx-auto px-4">
                <div className="flex items-center justify-between h-16">
                    {/* Logo */}
                    <NavLink to="/game" className="flex items-center gap-2">
                        <Target className={`w-8 h-8 ${darkMode ? theme.textColorDark : theme.textColor}`} />
                        <span className={`text-xl font-bold ${darkMode ? 'text-white' : 'text-gray-800'}`}>Darts</span>
                    </NavLink>

                    {/* Nav Links */}
                    <div className="flex items-center gap-1">
                        {navLinks.map(({ to, label, icon: Icon }) => (
                            <NavLink
                                key={to}
                                to={to}
                                className={({ isActive }) =>
                                    `flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-colors ${
                                        isActive
                                            ? `${theme.gradient} text-white`
                                            : darkMode
                                                ? 'text-gray-300 hover:bg-gray-800'
                                                : 'text-gray-600 hover:bg-gray-100'
                                    }`
                                }
                            >
                                <Icon className="w-5 h-5" />
                                <span className="hidden sm:inline">{label}</span>
                            </NavLink>
                        ))}
                    </div>

                    {/* Right side: Dark Mode + Theme + User */}
                    <div className="flex items-center gap-2">
                        {/* Dark Mode Toggle */}
                        <button
                            onClick={toggleDarkMode}
                            className={`p-2 rounded-lg transition-colors ${darkMode ? 'hover:bg-gray-800' : 'hover:bg-gray-100'}`}
                            title={darkMode ? 'Switch to light mode' : 'Switch to dark mode'}
                        >
                            {darkMode ? (
                                <Sun className="w-5 h-5 text-yellow-400" />
                            ) : (
                                <Moon className={`w-5 h-5 ${theme.textColor}`} />
                            )}
                        </button>

                        {/* Theme Switcher */}
                        <div className="relative" ref={themeMenuRef}>
                            <button
                                onClick={() => setShowThemeMenu(!showThemeMenu)}
                                className={`p-2 rounded-lg transition-colors ${darkMode ? 'hover:bg-gray-800' : 'hover:bg-gray-100'}`}
                                title="Change theme"
                            >
                                <Palette className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`} />
                            </button>
                            {showThemeMenu && (
                                <div className={`absolute right-0 mt-2 w-40 rounded-xl shadow-lg border py-2 ${darkMode ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-200'}`}>
                                    {Object.keys(THEMES).map((name) => (
                                        <button
                                            key={name}
                                            onClick={() => {
                                                setTheme(name);
                                                setShowThemeMenu(false);
                                            }}
                                            className={`w-full px-4 py-2 text-left flex items-center gap-3 ${
                                                darkMode ? 'hover:bg-gray-700' : 'hover:bg-gray-50'
                                            } ${themeName === name ? 'font-semibold' : ''} ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}
                                        >
                                            <span className={`w-4 h-4 rounded-full ${themeColors[name]}`} />
                                            <span className="capitalize">{name}</span>
                                        </button>
                                    ))}
                                </div>
                            )}
                        </div>

                        {/* User Menu */}
                        <div className="relative" ref={userMenuRef}>
                            <button
                                onClick={() => setShowUserMenu(!showUserMenu)}
                                className={`flex items-center gap-2 px-3 py-2 rounded-lg border-2 ${theme.borderColor} transition-colors ${darkMode ? 'hover:bg-gray-800' : 'hover:bg-gray-50'}`}
                            >
                                <User className={`w-5 h-5 ${darkMode ? theme.textColorDark : theme.textColor}`} />
                                <span className={`font-medium hidden sm:inline ${darkMode ? 'text-gray-200' : 'text-gray-700'}`}>
                                    {user?.username || 'User'}
                                </span>
                            </button>
                            {showUserMenu && (
                                <div className={`absolute right-0 mt-2 w-48 rounded-xl shadow-lg border py-2 ${darkMode ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-200'}`}>
                                    <NavLink
                                        to="/profile"
                                        onClick={() => setShowUserMenu(false)}
                                        className={`w-full px-4 py-2 text-left flex items-center gap-3 ${darkMode ? 'hover:bg-gray-700 text-gray-200' : 'hover:bg-gray-50 text-gray-700'}`}
                                    >
                                        <User className={`w-5 h-5 ${darkMode ? 'text-gray-400' : 'text-gray-500'}`} />
                                        <span>Profile</span>
                                    </NavLink>
                                    <NavLink
                                        to="/settings"
                                        onClick={() => setShowUserMenu(false)}
                                        className={`w-full px-4 py-2 text-left flex items-center gap-3 ${darkMode ? 'hover:bg-gray-700 text-gray-200' : 'hover:bg-gray-50 text-gray-700'}`}
                                    >
                                        <Settings className={`w-5 h-5 ${darkMode ? 'text-gray-400' : 'text-gray-500'}`} />
                                        <span>Settings</span>
                                    </NavLink>
                                    <hr className={`my-2 ${darkMode ? 'border-gray-700' : 'border-gray-200'}`} />
                                    <button
                                        onClick={() => {
                                            clearGame();
                                            logout();
                                            setShowUserMenu(false);
                                        }}
                                        className={`w-full px-4 py-2 text-left flex items-center gap-3 text-red-500 ${darkMode ? 'hover:bg-red-900/30' : 'hover:bg-red-50'}`}
                                    >
                                        <LogOut className="w-5 h-5" />
                                        <span>Logout</span>
                                    </button>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </nav>
    );
}