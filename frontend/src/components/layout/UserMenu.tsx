import { memo, useState, useRef, useCallback } from 'react';
import { NavLink } from 'react-router-dom';
import { User, Settings, LogOut } from '../ui/BoxIcon';
import { useAuth } from '../../contexts/AuthContext';
import { useGame } from '../../contexts/GameContext';
import { useClickOutside } from '../../hooks/useClickOutside';

export default memo(function UserMenu() {
    const { user, logout } = useAuth();
    const { clearAllGames } = useGame();
    const [open, setOpen] = useState(false);
    const menuRef = useRef<HTMLDivElement>(null);

    useClickOutside(menuRef, useCallback(() => setOpen(false), []));

    return (
        <div className="relative" ref={menuRef}>
            <button
                onClick={() => setOpen(!open)}
                className="flex items-center gap-2 p-1.5 sm:px-3 sm:py-2 rounded-lg border-2 border-theme-border transition-colors hover:bg-gray-50 dark:hover:bg-gray-800"
                aria-label="User menu"
                aria-expanded={open}
            >
                <User
                    className="w-4 h-4 sm:w-5 sm:h-5 text-theme-text"
                    aria-hidden="true"
                />
                <span className="font-medium hidden sm:inline text-gray-700 dark:text-gray-200">
                    {user?.username || 'User'}
                </span>
            </button>
            {open && (
                <div
                    className="absolute right-0 mt-2 w-48 rounded-xl shadow-lg border py-2 bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700"
                    role="menu"
                >
                    <NavLink
                        to="/profile"
                        onClick={() => setOpen(false)}
                        className="w-full px-4 py-2 text-left flex items-center gap-3 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-200"
                        role="menuitem"
                    >
                        <User className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                        <span>Profile</span>
                    </NavLink>
                    <NavLink
                        to="/settings"
                        onClick={() => setOpen(false)}
                        className="w-full px-4 py-2 text-left flex items-center gap-3 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-200"
                        role="menuitem"
                    >
                        <Settings className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                        <span>Settings</span>
                    </NavLink>
                    <hr className="my-2 border-gray-200 dark:border-gray-700" />
                    <button
                        onClick={() => {
                            clearAllGames();
                            logout();
                            setOpen(false);
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
    );
});
