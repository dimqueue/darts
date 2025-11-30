const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class RealApiClient {
    constructor() {
        this.baseURL = API_BASE_URL;
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

        const config = {
            ...options,
            headers,
        };

        try {
            const response = await fetch(url, config);

            if (response.status === 401) {
                localStorage.removeItem('token');
                localStorage.removeItem('user');
                window.location.href = '/auth';
                throw new Error('Session expired. Please log in again.');
            }

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            if (error.name === 'TypeError' && error.message.includes('fetch')) {
                throw new Error('Connection error. Please check your internet.');
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


    async getGlobalLeaderboard(params = { limit: 50, offset: 0 }) {
        const query = new URLSearchParams(params).toString();
        return this.request(`/api/leaderboard/global?${query}`, { method: 'GET' });
    }

    async getWeeklyLeaderboard(params = { limit: 50, offset: 0 }) {
        const query = new URLSearchParams(params).toString();
        return this.request(`/api/leaderboard/weekly?${query}`, { method: 'GET' });
    }

    async getMonthlyLeaderboard(params = { limit: 50, offset: 0 }) {
        const query = new URLSearchParams(params).toString();
        return this.request(`/api/leaderboard/monthly?${query}`, { method: 'GET' });
    }

    async getLanguageLeaderboard(language, params = { limit: 50, offset: 0 }) {
        const query = new URLSearchParams(params).toString();
        return this.request(`/api/leaderboard/language/${language}?${query}`, { method: 'GET' });
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