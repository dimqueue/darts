import { NavLink } from 'react-router-dom';
import { Trophy, User, BarChart3, Settings, LogOut, Palette } from 'lucide-react';
import { useAuth } from '../../contexts/AuthContext';
import { useTheme, THEMES } from '../../contexts/ThemeContext';
import { useState, useRef, useEffect } from 'react';

export default function Navbar() {
    const { user, logout } = useAuth();
    const { themeName, setTheme, theme } = useTheme();
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
        { to: '/game', label: 'Game', icon: Trophy },
        { to: '/leaderboard', label: 'Leaderboard', icon: BarChart3 },
    ];

    const themeColors = {
        purple: 'bg-purple-500',
        blue: 'bg-blue-500',
        green: 'bg-green-500',
    };

    return (
        <nav className="sticky top-0 z-50 bg-white shadow-md">
            <div className="max-w-6xl mx-auto px-4">
                <div className="flex items-center justify-between h-16">
                    {/* Logo */}
                    <NavLink to="/game" className="flex items-center gap-2">
                        <Trophy className={`w-8 h-8 ${theme.textColor}`} />
                        <span className="text-xl font-bold text-gray-800">Darts</span>
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
                                            : 'text-gray-600 hover:bg-gray-100'
                                    }`
                                }
                            >
                                <Icon className="w-5 h-5" />
                                <span className="hidden sm:inline">{label}</span>
                            </NavLink>
                        ))}
                    </div>

                    {/* Right side: Theme + User */}
                    <div className="flex items-center gap-2">
                        {/* Theme Switcher */}
                        <div className="relative" ref={themeMenuRef}>
                            <button
                                onClick={() => setShowThemeMenu(!showThemeMenu)}
                                className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
                                title="Change theme"
                            >
                                <Palette className={`w-5 h-5 ${theme.textColor}`} />
                            </button>
                            {showThemeMenu && (
                                <div className="absolute right-0 mt-2 w-40 bg-white rounded-xl shadow-lg border border-gray-200 py-2">
                                    {Object.keys(THEMES).map((name) => (
                                        <button
                                            key={name}
                                            onClick={() => {
                                                setTheme(name);
                                                setShowThemeMenu(false);
                                            }}
                                            className={`w-full px-4 py-2 text-left flex items-center gap-3 hover:bg-gray-50 ${
                                                themeName === name ? 'font-semibold' : ''
                                            }`}
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
                                className={`flex items-center gap-2 px-3 py-2 rounded-lg border-2 ${theme.borderColor} hover:bg-gray-50 transition-colors`}
                            >
                                <User className={`w-5 h-5 ${theme.textColor}`} />
                                <span className="font-medium text-gray-700 hidden sm:inline">
                                    {user?.username || 'User'}
                                </span>
                            </button>
                            {showUserMenu && (
                                <div className="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-200 py-2">
                                    <NavLink
                                        to="/profile"
                                        onClick={() => setShowUserMenu(false)}
                                        className="w-full px-4 py-2 text-left flex items-center gap-3 hover:bg-gray-50"
                                    >
                                        <User className="w-5 h-5 text-gray-500" />
                                        <span>Profile</span>
                                    </NavLink>
                                    <NavLink
                                        to="/settings"
                                        onClick={() => setShowUserMenu(false)}
                                        className="w-full px-4 py-2 text-left flex items-center gap-3 hover:bg-gray-50"
                                    >
                                        <Settings className="w-5 h-5 text-gray-500" />
                                        <span>Settings</span>
                                    </NavLink>
                                    <hr className="my-2" />
                                    <button
                                        onClick={() => {
                                            logout();
                                            setShowUserMenu(false);
                                        }}
                                        className="w-full px-4 py-2 text-left flex items-center gap-3 hover:bg-red-50 text-red-600"
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