import { createContext, useContext, useState, useEffect, useCallback, useMemo, type ReactNode } from 'react';
import { STORAGE_KEYS, THEME_OPTIONS, type ThemeName } from '../config/constants';
export type { ThemeName };

const THEME_NAMES = THEME_OPTIONS.map((opt) => opt.name);
const DEFAULT_THEME: ThemeName = 'purple';

interface ThemeContextValue {
    themeName: ThemeName;
    setTheme: (name: ThemeName) => void;
    darkMode: boolean;
    toggleDarkMode: () => void;
    setDarkMode: (value: boolean) => void;
}

const ThemeContext = createContext<ThemeContextValue | null>(null);

interface ThemeProviderProps {
    children: ReactNode;
}

function isValidTheme(name: string): name is ThemeName {
    return (THEME_NAMES as string[]).includes(name);
}

export function ThemeProvider({ children }: ThemeProviderProps) {
    const [themeName, setThemeName] = useState<ThemeName>(DEFAULT_THEME);
    const [darkMode, setDarkMode] = useState(false);

    // On mount: restore saved theme + dark mode, apply to DOM
    useEffect(() => {
        const savedTheme = localStorage.getItem(STORAGE_KEYS.THEME);
        if (savedTheme && isValidTheme(savedTheme)) {
            setThemeName(savedTheme);
            document.documentElement.dataset.theme = savedTheme;
        } else {
            document.documentElement.dataset.theme = DEFAULT_THEME;
        }

        const savedDarkMode = localStorage.getItem(STORAGE_KEYS.DARK_MODE);
        if (savedDarkMode !== null) {
            const isDark = savedDarkMode === 'true';
            setDarkMode(isDark);
            if (isDark) {
                document.documentElement.classList.add('dark');
            }
        } else {
            const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
            setDarkMode(prefersDark);
            if (prefersDark) {
                document.documentElement.classList.add('dark');
            }
        }
    }, []);

    // Sync dark class with state
    useEffect(() => {
        if (darkMode) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    }, [darkMode]);

    const applyDarkMode = useCallback((value: boolean) => {
        document.documentElement.classList.add('no-transitions');

        if (value) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }

        // Force reflow so class changes apply before transitions re-enable
        document.documentElement.offsetHeight;

        requestAnimationFrame(() => {
            document.documentElement.classList.remove('no-transitions');
        });
    }, []);

    const setTheme = useCallback((name: ThemeName) => {
        localStorage.setItem(STORAGE_KEYS.THEME, name);
        document.documentElement.dataset.theme = name;
        setThemeName(name);
    }, []);

    const toggleDarkMode = useCallback(() => {
        setDarkMode((prev) => {
            const newValue = !prev;
            localStorage.setItem(STORAGE_KEYS.DARK_MODE, String(newValue));
            applyDarkMode(newValue);
            return newValue;
        });
    }, [applyDarkMode]);

    const setDarkModeValue = useCallback((value: boolean) => {
        localStorage.setItem(STORAGE_KEYS.DARK_MODE, String(value));
        applyDarkMode(value);
        setDarkMode(value);
    }, [applyDarkMode]);

    const value = useMemo<ThemeContextValue>(
        () => ({
            themeName,
            setTheme,
            darkMode,
            toggleDarkMode,
            setDarkMode: setDarkModeValue,
        }),
        [themeName, setTheme, darkMode, toggleDarkMode, setDarkModeValue]
    );

    return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme(): ThemeContextValue {
    const context = useContext(ThemeContext);
    if (!context) {
        throw new Error('useTheme must be used within a ThemeProvider');
    }
    return context;
}

export default ThemeContext;
