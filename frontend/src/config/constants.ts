import type { Language } from '../types';

export const LANGUAGES: Language[] = [
    { code: 'en', name: 'English', flag: '\u{1F1EC}\u{1F1E7}' },
    { code: 'ua', name: 'Ukrainian', flag: '\u{1F1FA}\u{1F1E6}' },
];

export const LANGUAGES_WITH_ALL: Language[] = [
    { code: '', name: 'All Languages', flag: '\u{1F310}' },
    ...LANGUAGES,
];

export const DEFAULT_LANGUAGE = 'en';

export const PAGINATION = {
    DEFAULT_LIMIT: 50,
    MAX_LIMIT: 100,
    SHOW_PAGES: 5,
} as const;

export const TIMEOUTS = {
    FOCUS_DELAY: 100,
    FOCUS_DELAY_SHORT: 50,
    SUCCESS_MESSAGE: 3000,
    ERROR_MESSAGE: 5000,
    TOAST_INFO: 3000,
    TOAST_WARNING: 4000,
    DEBOUNCE: 300,
} as const;

export const GAME = {
    MIN_WORD_LENGTH: 2,
    MAX_WORD_LENGTH: 50,
    DISTANCE_THRESHOLDS: {
        VERY_CLOSE: 100,
        CLOSE: 500,
        MEDIUM: 1000,
        FAR: 5000,
    },
} as const;

export const PATTERNS = {
    WORD: /^[\p{L}\p{M}]+$/u,
    USERNAME: /^[a-zA-Z0-9_]+$/,
    COUNTRY_CODE: /^[A-Z]{2}$/,
} as const;

export const STORAGE_KEYS = {
    TOKEN: 'token',
    USER: 'user',
    THEME: 'theme',
    DARK_MODE: 'darkMode',
    GAME_CACHE_PREFIX: 'darts_game_cache_',
} as const;

export const ERROR_CODES = {
    UNAUTHORIZED: 'UNAUTHORIZED',
    FORBIDDEN: 'FORBIDDEN',
    NOT_FOUND: 'NOT_FOUND',
    VALIDATION_ERROR: 'VALIDATION_ERROR',
    GAME_NOT_FOUND: 'GAME_NOT_FOUND',
    GAME_NOT_ACTIVE: 'GAME_NOT_ACTIVE',
    WORD_NOT_FOUND: 'WORD_NOT_FOUND',
    WORD_ALREADY_USED: 'WORD_ALREADY_USED',
    USER_EXISTS: 'USER_EXISTS',
    CONNECTION_ERROR: 'CONNECTION_ERROR',
    INTERNAL_ERROR: 'INTERNAL_ERROR',
} as const;

export type ErrorCode = (typeof ERROR_CODES)[keyof typeof ERROR_CODES];

export const ERROR_MESSAGES: Record<string, string> = {
    [ERROR_CODES.UNAUTHORIZED]: 'Invalid credentials. Please try again.',
    [ERROR_CODES.FORBIDDEN]: "You don't have permission to do this.",
    [ERROR_CODES.NOT_FOUND]: 'The requested resource was not found.',
    [ERROR_CODES.GAME_NOT_FOUND]: 'Game not found.',
    [ERROR_CODES.GAME_NOT_ACTIVE]: 'This game is no longer active.',
    [ERROR_CODES.WORD_NOT_FOUND]: 'Word not found in vocabulary.',
    [ERROR_CODES.WORD_ALREADY_USED]: 'You already guessed this word.',
    [ERROR_CODES.USER_EXISTS]: 'Username already taken. Please choose another.',
    [ERROR_CODES.CONNECTION_ERROR]: 'Connection error. Please check your internet.',
    [ERROR_CODES.INTERNAL_ERROR]: 'Something went wrong. Please try again.',
    default: 'An unexpected error occurred.',
};

export interface AppError {
    code?: string;
    message?: string;
}

export function getErrorMessage(error: AppError | null | undefined): string {
    if (!error) return ERROR_MESSAGES.default;

    if (error.code && ERROR_MESSAGES[error.code]) {
        return ERROR_MESSAGES[error.code];
    }

    if (error.message) {
        return error.message;
    }

    return ERROR_MESSAGES.default;
}

export function getCountryFlag(countryCode: string | null | undefined): string | null {
    if (!countryCode || countryCode.length !== 2) return null;
    const codePoints = countryCode
        .toUpperCase()
        .split('')
        .map((char) => 127397 + char.charCodeAt(0));
    return String.fromCodePoint(...codePoints);
}
