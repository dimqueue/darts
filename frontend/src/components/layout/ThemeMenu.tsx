import { memo, useState, useRef, useCallback } from 'react';
import { Palette } from '../ui/BoxIcon';
import { useTheme } from '../../contexts/ThemeContext';
import { THEME_OPTIONS } from '../../config/constants';
import { useClickOutside } from '../../hooks/useClickOutside';

export default memo(function ThemeMenu() {
    const { themeName, setTheme } = useTheme();
    const [open, setOpen] = useState(false);
    const menuRef = useRef<HTMLDivElement>(null);

    useClickOutside(menuRef, useCallback(() => setOpen(false), []));

    return (
        <div className="relative" ref={menuRef}>
            <button
                onClick={() => setOpen(!open)}
                className="p-1.5 sm:p-2 rounded-lg transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
                aria-label="Change theme"
                aria-expanded={open}
            >
                <Palette className="w-4 h-4 sm:w-5 sm:h-5 text-theme-text" />
            </button>
            {open && (
                <div
                    className="absolute right-0 mt-2 w-40 rounded-xl shadow-lg border py-2 bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700"
                    role="menu"
                >
                    {THEME_OPTIONS.map((opt) => (
                        <button
                            key={opt.name}
                            onClick={() => {
                                setTheme(opt.name);
                                setOpen(false);
                            }}
                            className={`w-full px-4 py-2 text-left flex items-center gap-3 hover:bg-gray-50 dark:hover:bg-gray-700 ${themeName === opt.name ? 'font-semibold' : ''} text-gray-700 dark:text-gray-200`}
                            role="menuitem"
                        >
                            <span className={`w-4 h-4 rounded-full ${opt.swatch}`} />
                            <span className="capitalize">{opt.name}</span>
                        </button>
                    ))}
                </div>
            )}
        </div>
    );
});
