import type { ThemeName } from '../config/constants';
import type { User } from './index';

export interface ApiResponse<T> {
    data?: T;
    message?: string;
}

export interface AuthResponse {
    token: string;
    user_id: number;
}

export interface SignUpResponse {
    message: string;
}

export interface ProfileData extends User {
    total_games?: number;
    total_wins?: number;
    total_losses?: number;
    current_win_streak?: number;
    best_win_streak?: number;
    average_guesses?: number;
    total_score?: number;
}

export interface ProfileUpdateData {
    bio?: string;
    country_code?: string;
}

export type EditFormData = Required<ProfileUpdateData>;

export interface SettingsData {
    preferred_language: string;
    theme: ThemeName;
    sound_enabled: boolean;
    email_notifications: boolean;
    show_profile_public: boolean;
    show_stats_public: boolean;
}

export interface StatisticsData {
    total_games: number;
    total_wins: number;
    total_losses?: number;
    current_win_streak: number;
    best_win_streak: number;
    total_guesses: number;
    average_guesses: number;
    fastest_win_seconds?: number;
    fewest_guesses_win?: number;
    total_score: number;
}

export interface LanguageStatsData {
    language: string;
    games_played: number;
    games_won: number;
    average_guesses: number;
    total_score?: number;
}

export interface LeaderboardEntry {
    rank: number;
    user_id: number;
    username: string;
    name?: string;
    avatar_url?: string | null;
    country_code?: string;
    total_score: number;
    total_wins: number;
    total_games?: number;
    best_win_streak: number;
    average_guesses: number;
    win_rate: number | string;
}

export interface LeaderboardResponse {
    leaderboard_type: string;
    language?: string | null;
    users: LeaderboardEntry[];
    total: number;
    current_user_rank?: number | null;
}

export interface RankData {
    global_rank: number;
    daily_rank: number;
    weekly_rank: number;
    monthly_rank: number;
}

export interface PublicProfileData {
    username: string;
    name?: string;
    avatar_url?: string | null;
    country_code?: string;
    total_games: number;
    total_wins: number;
    win_rate: number;
    best_win_streak: number;
}

export class ApiError extends Error {
    code: string;
    status: number;
    requestId: string | null;

    constructor(message: string, code: string, status: number, requestId: string | null = null) {
        super(message);
        this.name = 'ApiError';
        this.code = code;
        this.status = status;
        this.requestId = requestId;
    }
}
