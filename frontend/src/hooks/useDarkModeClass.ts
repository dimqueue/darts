import { useTheme } from '../contexts/ThemeContext';

type ClassMapping = {
    light: string;
    dark: string;
};

/**
 * Hook for generating dark mode aware class names
 * @param classes - Object with light and dark class strings, or a function that returns classes
 * @returns The appropriate class string based on current theme
 */
export function useDarkModeClass(classes: ClassMapping): string {
    const { darkMode } = useTheme();
    return darkMode ? classes.dark : classes.light;
}

/**
 * Utility function to create conditional dark mode classes without hook
 * @param darkMode - Current dark mode state
 * @param classes - Object with light and dark class strings
 * @returns The appropriate class string
 */
export function getDarkModeClass(darkMode: boolean, classes: ClassMapping): string {
    return darkMode ? classes.dark : classes.light;
}

/**
 * Common dark mode class presets
 */
export const darkModePresets = {
    card: {
        light: 'bg-white text-gray-800 border-gray-200',
        dark: 'bg-gray-800 text-gray-100 border-gray-700',
    },
    input: {
        light: 'bg-white text-gray-800 border-gray-300 placeholder-gray-400',
        dark: 'bg-gray-700 text-white border-gray-600 placeholder-gray-400',
    },
    dropdown: {
        light: 'bg-white border-gray-200',
        dark: 'bg-gray-800 border-gray-700',
    },
    hover: {
        light: 'hover:bg-gray-100',
        dark: 'hover:bg-gray-700',
    },
    text: {
        light: 'text-gray-800',
        dark: 'text-gray-100',
    },
    textMuted: {
        light: 'text-gray-600',
        dark: 'text-gray-400',
    },
} as const;
