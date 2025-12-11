import config from '../config/env';

// Custom error class with code and request_id
export class ApiError extends Error {
    constructor(message, code, status, requestId = null) {
        super(message);
        this.name = 'ApiError';
        this.code = code;
        this.status = status;
        this.requestId = requestId;
    }
}

class RealApiClient {
    constructor() {
        this.baseURL = config.apiUrl;
        console.log(`API Mode: REAL (${this.baseURL})`);
    }

    getToken() {
        return localStorage.getItem('token');
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        };

        const token = this.getToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const fetchConfig = {
            ...options,
            headers,
        };

        try {
            const response = await fetch(url, fetchConfig);
            const data = await response.json().catch(() => ({}));

            const isAuthEndpoint = endpoint.startsWith('/auth/');
            if (response.status === 401 && !isAuthEndpoint) {
                localStorage.removeItem('token');
                localStorage.removeItem('user');
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

            return data;
        } catch (error) {
            if (error instanceof ApiError) {
                throw error;
            }
            if (error.name === 'TypeError' && error.message.includes('fetch')) {
                throw new ApiError(
                    'Connection error. Please check your internet.',
                    'CONNECTION_ERROR',
                    0
                );
            }
            throw error;
        }
    }

    async signUp(username, password, name = '') {
        return this.request('/auth/sign-up', {
            method: 'POST',
            body: JSON.stringify({
                username,
                password,
                name: name || username,
            }),
        });
    }

    async signIn(username, password) {
        return this.request('/auth/sign-in', {
            method: 'POST',
            body: JSON.stringify({
                username,
                password,
            }),
        });
    }

    async createGame(language) {
        return this.request('/api/games', {
            method: 'POST',
            body: JSON.stringify({ language }),
        });
    }

    async getAllGames() {
        return this.request('/api/games', { method: 'GET' });
    }

    async getActiveGame() {
        return this.request('/api/games/active', { method: 'GET' });
    }

    async getGameById(gameId) {
        return this.request(`/api/games/${gameId}`, { method: 'GET' });
    }

    async updateGame(gameId, data) {
        return this.request(`/api/games/${gameId}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteGame(gameId) {
        return this.request(`/api/games/${gameId}`, { method: 'DELETE' });
    }

    async abandonGame(gameId) {
        return this.request(`/api/games/${gameId}/abandon`, { method: 'POST' });
    }

    async createGuess(gameId, guess) {
        return this.request(`/api/games/${gameId}/guesses`, {
            method: 'POST',
            body: JSON.stringify({ guess }),
        });
    }

    async getAllGuessByGame(gameId) {
        return this.request(`/api/games/${gameId}/guesses`, { method: 'GET' });
    }

    async getGuessById(gameId, guessId) {
        return this.request(`/api/games/${gameId}/guesses/${guessId}`, { method: 'GET' });
    }

    async getMyProfile() {
        return this.request('/api/profile', { method: 'GET' });
    }

    async updateMyProfile(data) {
        return this.request('/api/profile', {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async getMySettings() {
        return this.request('/api/profile/settings', { method: 'GET' });
    }

    async updateMySettings(data) {
        return this.request('/api/profile/settings', {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async getMyStatistics() {
        return this.request('/api/profile/statistics', { method: 'GET' });
    }

    async getMyLanguageStats() {
        return this.request('/api/profile/statistics/languages', { method: 'GET' });
    }

    async getLeaderboard(type = 'global', params = { limit: 50, offset: 0 }, language = null) {
        const queryParams = { type, ...params };
        if (language) {
            queryParams.language = language;
        }
        const query = new URLSearchParams(queryParams).toString();
        return this.request(`/api/leaderboard?${query}`, { method: 'GET' });
    }

    async getMyRank() {
        return this.request('/api/leaderboard/my-rank', { method: 'GET' });
    }

    async getPublicProfile(username) {
        return this.request(`/public/profile/${username}`, { method: 'GET' });
    }

    clearAllData() {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        localStorage.removeItem('theme');
    }
}

export default RealApiClient;
