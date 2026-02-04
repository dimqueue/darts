import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import type { Theme } from '../types';

export const THEMES: Record<string, Theme> = {
    purple: {
        name: 'purple',
        primary: '#8B5CF6',
        secondary: '#A78BFA',
        gradient: 'bg-gradient-purple',
        bgGradient: 'from-violet-100 via-purple-100 to-fuchsia-100',
        bgGradientDark: 'from-gray-900 via-violet-950 to-gray-900',
        meshBg: 'bg-mesh-purple',
        focusBorder: 'focus:border-violet-500',
        textColor: 'text-violet-600',
        textColorDark: 'text-violet-400',
        borderColor: 'border-violet-500',
        hoverBg: 'hover:bg-violet-50',
        hoverBgDark: 'hover:bg-violet-900/30',
        lightBg: 'bg-violet-50',
        lightBgDark: 'bg-violet-900/20',
    },
    blue: {
        name: 'blue',
        primary: '#3B82F6',
        secondary: '#60A5FA',
        gradient: 'bg-gradient-blue',
        bgGradient: 'from-blue-100 via-sky-100 to-cyan-100',
        bgGradientDark: 'from-gray-900 via-blue-950 to-gray-900',
        meshBg: 'bg-mesh-blue',
        focusBorder: 'focus:border-blue-500',
        textColor: 'text-blue-600',
        textColorDark: 'text-blue-400',
        borderColor: 'border-blue-500',
        hoverBg: 'hover:bg-blue-50',
        hoverBgDark: 'hover:bg-blue-900/30',
        lightBg: 'bg-blue-50',
        lightBgDark: 'bg-blue-900/20',
    },
    green: {
        name: 'green',
        primary: '#10B981',
        secondary: '#34D399',
        gradient: 'bg-gradient-green',
        bgGradient: 'from-emerald-50 via-green-50 to-teal-50',
        bgGradientDark: 'from-gray-900 via-emerald-950 to-gray-900',
        meshBg: 'bg-mesh-green',
        focusBorder: 'focus:border-emerald-500',
        textColor: 'text-emerald-600',
        textColorDark: 'text-emerald-400',
        borderColor: 'border-emerald-500',
        hoverBg: 'hover:bg-emerald-50',
        hoverBgDark: 'hover:bg-emerald-900/30',
        lightBg: 'bg-emerald-50',
        lightBgDark: 'bg-emerald-900/20',
    },
};

const DEFAULT_THEME = 'purple';

interface ThemeContextValue {
    themeName: string;
    theme: Theme;
    setTheme: (name: string) => void;
    themes: typeof THEMES;
    darkMode: boolean;
    toggleDarkMode: () => void;
    setDarkMode: (value: boolean) => void;
}

const ThemeContext = createContext<ThemeContextValue | null>(null);

interface ThemeProviderProps {
    children: ReactNode;
}

export function ThemeProvider({ children }: ThemeProviderProps) {
    const [themeName, setThemeName] = useState(DEFAULT_THEME);
    const [darkMode, setDarkMode] = useState(false);

    useEffect(() => {
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme && THEMES[savedTheme]) {
            setThemeName(savedTheme);
        }

        const savedDarkMode = localStorage.getItem('darkMode');
        if (savedDarkMode !== null) {
            setDarkMode(savedDarkMode === 'true');
        } else {
            const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
            setDarkMode(prefersDark);
        }
    }, []);

    useEffect(() => {
        if (darkMode) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    }, [darkMode]);


    const setTheme = (name: string) => {
        if (THEMES[name]) {
            localStorage.setItem('theme', name);
            setThemeName(name);
        }
    };

    const toggleDarkMode = () => {
        const newValue = !darkMode;
        localStorage.setItem('darkMode', String(newValue));
        setDarkMode(newValue);
    };

    const setDarkModeValue = (value: boolean) => {
        localStorage.setItem('darkMode', String(value));
        setDarkMode(value);
    };

    const theme = THEMES[themeName];

    const value: ThemeContextValue = {
        themeName,
        theme,
        setTheme,
        themes: THEMES,
        darkMode,
        toggleDarkMode,
        setDarkMode: setDarkModeValue,
    };

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
