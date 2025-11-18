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

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
                throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('API request failed:', error);
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
            body: JSON.stringify({
                language: language,
            }),
        });
    }

    async getAllGames() {
        return this.request('/api/games', {
            method: 'GET',
        });
    }

    async getGameById(gameId) {
        return this.request(`/api/games/${gameId}`, {
            method: 'GET',
        });
    }

    async updateGame(gameId, data) {
        return this.request(`/api/games/${gameId}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteGame(gameId) {
        return this.request(`/api/games/${gameId}`, {
            method: 'DELETE',
        });
    }

    async createGuess(gameId, guess) {
        return this.request(`/api/games/${gameId}/guesses`, {
            method: 'POST',
            body: JSON.stringify({
                guess,
            }),
        });
    }

    async getAllGuessByGame(gameId) {
        return this.request(`/api/games/${gameId}/guesses`, {
            method: 'GET',
        });
    }

    async getGuessById(gameId, guessId) {
        return this.request(`/api/games/${gameId}/guesses/${guessId}`, {
            method: 'GET',
        });
    }

    async clearAllData() {
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        localStorage.removeItem('currentGameId');
    }
}

export default RealApiClient;