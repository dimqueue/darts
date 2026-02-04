import { NavLink } from 'react-router-dom';
import {
    Target,
    Crosshair,
    Layers,
    Crown,
    User,
    Settings,
    LogOut,
    Palette,
    Moon,
    Sun,
} from '../ui/BoxIcon';
import { useAuth } from '../../contexts/AuthContext';
import { useGame } from '../../contexts/GameContext';
import { useTheme, THEMES } from '../../contexts/ThemeContext';
import { useState, useRef, useEffect } from 'react';

const navLinks = [
    { to: '/game', label: 'Play', icon: Crosshair },
    { to: '/modes', label: 'Modes', icon: Layers },
    { to: '/leaderboard', label: 'Leaderboard', icon: Crown },
    // { to: '/new', label: 'New', icon: Circle }, // TODO: Placeholder - replace icon and route
];

const themeColors: Record<string, string> = {
    purple: 'bg-purple-500',
    blue: 'bg-blue-500',
    green: 'bg-green-500',
};

export default function Navbar() {
    const { user, logout } = useAuth();
    const { clearGame } = useGame();
    const { themeName, setTheme, theme, darkMode, toggleDarkMode } = useTheme();
    const [showThemeMenu, setShowThemeMenu] = useState(false);
    const [showUserMenu, setShowUserMenu] = useState(false);
    const themeMenuRef = useRef<HTMLDivElement>(null);
    const userMenuRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        function handleClickOutside(event: MouseEvent) {
            if (themeMenuRef.current && !themeMenuRef.current.contains(event.target as Node)) {
                setShowThemeMenu(false);
            }
            if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
                setShowUserMenu(false);
            }
        }
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    return (
        <nav
            className="sticky top-0 z-50 shadow-md bg-white dark:bg-gray-900"
            aria-label="Main navigation"
        >
            <div className="max-w-6xl mx-auto px-2 sm:px-4">
                <div className="grid grid-cols-3 items-center h-14 sm:h-16">
                    {/* Left: Logo */}
                    <NavLink to="/game" className="flex items-center gap-1 sm:gap-2 justify-self-start">
                        <Target
                            className={`w-6 h-6 sm:w-8 sm:h-8 ${darkMode ? theme.textColorDark : theme.textColor}`}
                            aria-hidden="true"
                        />
                        <span className="text-lg sm:text-xl font-bold hidden xs:inline text-gray-800 dark:text-white">
                            Darts
                        </span>
                    </NavLink>

                    {/* Center: Nav links */}
                    <div className="flex items-center justify-center gap-0.5 sm:gap-1">
                        {navLinks.map((link) => {
                            const LinkIcon = link.icon;
                            return (
                                <NavLink
                                    key={link.to}
                                    to={link.to}
                                    className={({ isActive }) =>
                                        `flex items-center gap-1 sm:gap-2 px-2 sm:px-4 py-1.5 sm:py-2 rounded-lg font-medium transition-colors ${
                                            isActive
                                                ? `${theme.gradient} text-white`
                                                : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800'
                                        }`
                                    }
                                >
                                    <LinkIcon className="w-4 h-4 sm:w-5 sm:h-5" aria-hidden="true" />
                                    <span className="hidden sm:inline">{link.label}</span>
                                </NavLink>
                            );
                        })}
                    </div>

                    {/* Right: Actions */}
                    <div className="flex items-center gap-0.5 sm:gap-2 justify-self-end">
                        <button
                            onClick={toggleDarkMode}
                            className="p-1.5 sm:p-2 rounded-lg transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
                            aria-label={darkMode ? 'Switch to light mode' : 'Switch to dark mode'}
                        >
                            {darkMode ? (
                                <Sun className="w-4 h-4 sm:w-5 sm:h-5 text-yellow-400" />
                            ) : (
                                <Moon className={`w-4 h-4 sm:w-5 sm:h-5 ${theme.textColor}`} />
                            )}
                        </button>

                        <div className="relative" ref={themeMenuRef}>
                            <button
                                onClick={() => setShowThemeMenu(!showThemeMenu)}
                                className="p-1.5 sm:p-2 rounded-lg transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
                                aria-label="Change theme"
                                aria-expanded={showThemeMenu}
                            >
                                <Palette
                                    className={`w-4 h-4 sm:w-5 sm:h-5 ${darkMode ? theme.textColorDark : theme.textColor}`}
                                />
                            </button>
                            {showThemeMenu && (
                                <div
                                    className="absolute right-0 mt-2 w-40 rounded-xl shadow-lg border py-2 bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700"
                                    role="menu"
                                >
                                    {Object.keys(THEMES).map((name) => (
                                        <button
                                            key={name}
                                            onClick={() => {
                                                setTheme(name);
                                                setShowThemeMenu(false);
                                            }}
                                            className={`w-full px-4 py-2 text-left flex items-center gap-3 hover:bg-gray-50 dark:hover:bg-gray-700 ${themeName === name ? 'font-semibold' : ''} text-gray-700 dark:text-gray-200`}
                                            role="menuitem"
                                        >
                                            <span
                                                className={`w-4 h-4 rounded-full ${themeColors[name]}`}
                                            />
                                            <span className="capitalize">{name}</span>
                                        </button>
                                    ))}
                                </div>
                            )}
                        </div>

                        <div className="relative" ref={userMenuRef}>
                            <button
                                onClick={() => setShowUserMenu(!showUserMenu)}
                                className={`flex items-center gap-2 p-1.5 sm:px-3 sm:py-2 rounded-lg border-2 ${theme.borderColor} transition-colors hover:bg-gray-50 dark:hover:bg-gray-800`}
                                aria-label="User menu"
                                aria-expanded={showUserMenu}
                            >
                                <User
                                    className={`w-4 h-4 sm:w-5 sm:h-5 ${darkMode ? theme.textColorDark : theme.textColor}`}
                                    aria-hidden="true"
                                />
                                <span className="font-medium hidden sm:inline text-gray-700 dark:text-gray-200">
                                    {user?.username || 'User'}
                                </span>
                            </button>
                            {showUserMenu && (
                                <div
                                    className="absolute right-0 mt-2 w-48 rounded-xl shadow-lg border py-2 bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700"
                                    role="menu"
                                >
                                    <NavLink
                                        to="/profile"
                                        onClick={() => setShowUserMenu(false)}
                                        className="w-full px-4 py-2 text-left flex items-center gap-3 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-200"
                                        role="menuitem"
                                    >
                                        <User className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                                        <span>Profile</span>
                                    </NavLink>
                                    <NavLink
                                        to="/settings"
                                        onClick={() => setShowUserMenu(false)}
                                        className="w-full px-4 py-2 text-left flex items-center gap-3 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-200"
                                        role="menuitem"
                                    >
                                        <Settings className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                                        <span>Settings</span>
                                    </NavLink>
                                    <hr className="my-2 border-gray-200 dark:border-gray-700" />
                                    <button
                                        onClick={() => {
                                            clearGame();
                                            logout();
                                            setShowUserMenu(false);
                                        }}
                                        className="w-full px-4 py-2 text-left flex items-center gap-3 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/30"
                                        role="menuitem"
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
