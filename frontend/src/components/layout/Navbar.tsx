import { NavLink } from 'react-router-dom';
import { Target, Home, Crown, Moon, Sun } from '../ui/BoxIcon';
import { useTheme } from '../../contexts/ThemeContext';
import ThemeMenu from './ThemeMenu';
import UserMenu from './UserMenu';

const navLinks = [
    { to: '/home', label: 'Home', icon: Home },
    { to: '/leaderboard', label: 'Leaderboard', icon: Crown },
] as const;

export default function Navbar() {
    const { darkMode, toggleDarkMode } = useTheme();

    return (
        <nav
            className="shrink-0 z-50 shadow-md bg-white dark:bg-gray-900"
            aria-label="Main navigation"
        >
            <div className="max-w-6xl mx-auto px-2 sm:px-4">
                <div className="grid grid-cols-3 items-center h-14 sm:h-16">
                    {/* Left: Logo */}
                    <NavLink to="/home" className="flex items-center gap-1 sm:gap-2 justify-self-start">
                        <Target
                            className="w-6 h-6 sm:w-8 sm:h-8 text-theme-text"
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
                                                ? 'bg-theme-gradient text-white'
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
                                <Moon className="w-4 h-4 sm:w-5 sm:h-5 text-gray-600 dark:text-gray-300" />
                            )}
                        </button>

                        <ThemeMenu />
                        <UserMenu />
                    </div>
                </div>
            </div>
        </nav>
    );
}
