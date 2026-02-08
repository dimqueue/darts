import config from '../config/env';
import { STORAGE_KEYS } from '../config/constants';
import {
    ApiError,
    type AuthResponse,
    type SignUpResponse,
    type ProfileData,
    type ProfileUpdateData,
    type SettingsData,
    type StatisticsData,
    type LanguageStatsData,
    type LeaderboardResponse,
    type RankData,
    type PublicProfileData,
    type ApiResponse,
} from '../types/api';
import type { Game, Guess, GuessResult } from '../types/game';

interface RequestOptions extends RequestInit {
    headers?: Record<string, string>;
}

class RealApiClient {
    private baseURL: string;

    constructor() {
        this.baseURL = config.apiUrl;
        console.log(`API Mode: REAL (${this.baseURL})`);
    }

    private getToken(): string | null {
        return localStorage.getItem(STORAGE_KEYS.TOKEN);
    }

    private async request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
        const url = `${this.baseURL}${endpoint}`;
        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
            ...options.headers,
        };

        const token = this.getToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const fetchConfig: RequestInit = {
            ...options,
            headers,
        };

        try {
            const response = await fetch(url, fetchConfig);
            const data = await response.json().catch(() => ({}));

            const isAuthEndpoint = endpoint.startsWith('/auth/');
            if (response.status === 401 && !isAuthEndpoint) {
                localStorage.removeItem(STORAGE_KEYS.TOKEN);
                localStorage.removeItem(STORAGE_KEYS.USER);
                window.location.href = '/auth';
                throw new ApiError('Session expired. Please log in again.', 'UNAUTHORIZED', 401);
            }

            if (!response.ok) {
                throw new ApiError(
                    data.message || `HTTP error! status: ${response.status}`,
                    data.code || 'UNKNOWN_ERROR',
                    response.status,
                    data.request_id
                );
            }

            return data as T;
        } catch (error) {
            if (error instanceof ApiError) {
                throw error;
            }
            if (error instanceof TypeError && error.message.includes('fetch')) {
                throw new ApiError(
                    'Connection error. Please check your internet.',
                    'CONNECTION_ERROR',
                    0
                );
            }
            throw error;
        }
    }

    async signUp(username: string, password: string, name: string = ''): Promise<SignUpResponse> {
        return this.request('/auth/sign-up', {
            method: 'POST',
            body: JSON.stringify({
                username,
                password,
                name: name || username,
            }),
        });
    }

    async signIn(username: string, password: string): Promise<AuthResponse> {
        return this.request('/auth/sign-in', {
            method: 'POST',
            body: JSON.stringify({
                username,
                password,
            }),
        });
    }

    async createGame(language: string): Promise<Game> {
        return this.request('/api/games', {
            method: 'POST',
            body: JSON.stringify({ language }),
        });
    }

    async getAllGames(): Promise<ApiResponse<Game[]>> {
        return this.request('/api/games', { method: 'GET' });
    }

    async getActiveGame(): Promise<ApiResponse<Game | null>> {
        return this.request('/api/games/active', { method: 'GET' });
    }

    async getGameById(gameId: number): Promise<Game> {
        return this.request(`/api/games/${gameId}`, { method: 'GET' });
    }

    async updateGame(gameId: number, data: Partial<Game>): Promise<Game> {
        return this.request(`/api/games/${gameId}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteGame(gameId: number): Promise<{ message: string }> {
        return this.request(`/api/games/${gameId}`, { method: 'DELETE' });
    }

    async abandonGame(gameId: number): Promise<{ message: string }> {
        return this.request(`/api/games/${gameId}/abandon`, { method: 'POST' });
    }

    async createGuess(gameId: number, guess: string): Promise<GuessResult> {
        return this.request(`/api/games/${gameId}/guesses`, {
            method: 'POST',
            body: JSON.stringify({ guess }),
        });
    }

    async getAllGuessByGame(gameId: number): Promise<ApiResponse<Guess[]>> {
        return this.request(`/api/games/${gameId}/guesses`, { method: 'GET' });
    }

    async getGuessById(gameId: number, guessId: number): Promise<Guess> {
        return this.request(`/api/games/${gameId}/guesses/${guessId}`, { method: 'GET' });
    }

    async getMyProfile(): Promise<ProfileData> {
        return this.request('/api/profile', { method: 'GET' });
    }

    async updateMyProfile(data: ProfileUpdateData): Promise<{ message: string }> {
        return this.request('/api/profile', {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async getMySettings(): Promise<SettingsData> {
        return this.request('/api/profile/settings', { method: 'GET' });
    }

    async updateMySettings(data: Partial<SettingsData>): Promise<{ message: string }> {
        return this.request('/api/profile/settings', {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async getMyStatistics(): Promise<StatisticsData> {
        return this.request('/api/profile/statistics', { method: 'GET' });
    }

    async getMyLanguageStats(): Promise<ApiResponse<LanguageStatsData[]>> {
        return this.request('/api/profile/statistics/languages', { method: 'GET' });
    }

    async getLeaderboard(
        type: string = 'global',
        params: { limit?: number; offset?: number } = { limit: 50, offset: 0 },
        language: string | null = null
    ): Promise<LeaderboardResponse> {
        const queryParams: Record<string, string | number> = { type, ...params };
        if (language) {
            queryParams.language = language;
        }
        const query = new URLSearchParams(
            Object.entries(queryParams).map(([k, v]) => [k, String(v)])
        ).toString();
        return this.request(`/api/leaderboard?${query}`, { method: 'GET' });
    }

    async getMyRank(): Promise<RankData> {
        return this.request('/api/leaderboard/my-rank', { method: 'GET' });
    }

    async getPublicProfile(username: string): Promise<PublicProfileData> {
        return this.request(`/public/profile/${username}`, { method: 'GET' });
    }

    clearAllData(): void {
        localStorage.removeItem(STORAGE_KEYS.TOKEN);
        localStorage.removeItem(STORAGE_KEYS.USER);
        localStorage.removeItem(STORAGE_KEYS.THEME);
    }
}

export default RealApiClient;
