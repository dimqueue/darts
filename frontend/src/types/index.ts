export * from './api';
export * from './game';

export interface User {
    id: number;
    username: string;
    name?: string;
    email?: string;
    avatar_url?: string | null;
    bio?: string;
    country_code?: string;
    created_at?: string;
}

export interface Theme {
    name: string;
    primary: string;
    secondary: string;
    gradient: string;
    bgGradient: string;
    bgGradientDark: string;
    meshBg: string;
    focusBorder: string;
    textColor: string;
    textColorDark: string;
    borderColor: string;
    hoverBg: string;
    hoverBgDark: string;
    lightBg: string;
    lightBgDark: string;
}

export interface ToastItem {
    id: number;
    message: string;
    type: 'success' | 'error' | 'info' | 'warning';
    duration: number;
}

export interface Language {
    code: string;
    name: string;
    flag: string;
}
