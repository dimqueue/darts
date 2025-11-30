import { createContext, useContext, useState, useEffect } from 'react';

const ThemeContext = createContext(null);

export const THEMES = {
    purple: {
        name: 'purple',
        primary: '#A855F7',
        secondary: '#EC4899',
        gradient: 'bg-gradient-purple',
        bgGradient: 'from-purple-50 to-pink-50',
        focusBorder: 'focus:border-purple-500',
        textColor: 'text-purple-600',
        borderColor: 'border-purple-500',
        hoverBg: 'hover:bg-purple-50',
    },
    blue: {
        name: 'blue',
        primary: '#3B82F6',
        secondary: '#06B6D4',
        gradient: 'bg-gradient-blue',
        bgGradient: 'from-blue-50 to-cyan-50',
        focusBorder: 'focus:border-blue-500',
        textColor: 'text-blue-600',
        borderColor: 'border-blue-500',
        hoverBg: 'hover:bg-blue-50',
    },
    green: {
        name: 'green',
        primary: '#10B981',
        secondary: '#3B82F6',
        gradient: 'bg-gradient-green',
        bgGradient: 'from-green-50 to-blue-50',
        focusBorder: 'focus:border-green-500',
        textColor: 'text-green-600',
        borderColor: 'border-green-500',
        hoverBg: 'hover:bg-green-50',
    },
};

const DEFAULT_THEME = 'purple';

export function ThemeProvider({ children }) {
    const [themeName, setThemeName] = useState(DEFAULT_THEME);

    useEffect(() => {
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme && THEMES[savedTheme]) {
            setThemeName(savedTheme);
        }
    }, []);

    const setTheme = (name) => {
        if (THEMES[name]) {
            localStorage.setItem('theme', name);
            setThemeName(name);
        }
    };

    const theme = THEMES[themeName];

    const value = {
        themeName,
        theme,
        setTheme,
        themes: THEMES,
    };

    return (
        <ThemeContext.Provider value={value}>
            {children}
        </ThemeContext.Provider>
    );
}

export function useTheme() {
    const context = useContext(ThemeContext);
    if (!context) {
        throw new Error('useTheme must be used within a ThemeProvider');
    }
    return context;
}

export default ThemeContext;