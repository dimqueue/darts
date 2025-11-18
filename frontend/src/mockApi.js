// Mock API for testing without backend

class MockApiClient {
    constructor() {
        this.users = JSON.parse(localStorage.getItem('mock_users') || '[]');
        this.games = JSON.parse(localStorage.getItem('mock_games') || '[]');
        this.guesses = JSON.parse(localStorage.getItem('mock_guesses') || '[]');
        this.currentUserId = null;
        this.nextGameId = this.games.length + 1;
        this.nextGuessId = this.guesses.length + 1;
    }

    // Helper to simulate network delay
    async delay(ms = 100) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    // Helper to save data
    save() {
        localStorage.setItem('mock_users', JSON.stringify(this.users));
        localStorage.setItem('mock_games', JSON.stringify(this.games));
        localStorage.setItem('mock_guesses', JSON.stringify(this.guesses));
    }

    // Helper to get token
    getToken() {
        return localStorage.getItem('token');
    }

    // Helper to parse token and get userId
    getUserIdFromToken() {
        const token = this.getToken();
        if (!token) return null;
        try {
            // Mock token format: "mock_token_userId"
            return parseInt(token.split('_')[2]);
        } catch {
            return null;
        }
    }

    // Auth endpoints
    async signUp(username, password, name = '') {
        await this.delay();

        // Check if user exists
        if (this.users.find(u => u.username === username)) {
            throw new Error('User already exists');
        }

        const userId = this.users.length + 1;
        const user = {
            id: userId,
            username,
            password,
            name: name || username,
        };

        this.users.push(user);
        this.save();

        const token = `mock_token_${userId}`;
        return { token };
    }

    async signIn(username, password) {
        await this.delay();

        const user = this.users.find(u => u.username === username && u.password === password);
        if (!user) {
            throw new Error('Invalid credentials');
        }

        const token = `mock_token_${user.id}`;
        return { token };
    }

    // Game endpoints
    async createGame(language) {
        await this.delay();

        const userId = this.getUserIdFromToken();
        if (!userId) throw new Error('Unauthorized');

        const game = {
            id: this.nextGameId++,
            user_id: userId,
            language: language || 'en',
            status: 'in_progress',
            word_id: 1,
            started_at: new Date().toISOString(),
            ended_at: null,
        };

        this.games.push(game);
        this.save();

        return game;
    }

    async getAllGames() {
        await this.delay();

        const userId = this.getUserIdFromToken();
        if (!userId) throw new Error('Unauthorized');

        return this.games.filter(g => g.user_id === userId);
    }

    async getGameById(gameId) {
        await this.delay();

        const userId = this.getUserIdFromToken();
        if (!userId) throw new Error('Unauthorized');

        const game = this.games.find(g => g.id === gameId);
        if (!game) throw new Error('Game not found');
        if (game.user_id !== userId) throw new Error('Unauthorized');

        return game;
    }

    async updateGame(gameId, data) {
        await this.delay();

        const userId = this.getUserIdFromToken();
        if (!userId) throw new Error('Unauthorized');

        const gameIndex = this.games.findIndex(g => g.id === gameId);
        if (gameIndex === -1) throw new Error('Game not found');
        if (this.games[gameIndex].user_id !== userId) throw new Error('Unauthorized');

        this.games[gameIndex] = { ...this.games[gameIndex], ...data };
        this.save();

        return this.games[gameIndex];
    }

    async deleteGame(gameId) {
        await this.delay();

        const userId = this.getUserIdFromToken();
        if (!userId) throw new Error('Unauthorized');

        const gameIndex = this.games.findIndex(g => g.id === gameId);
        if (gameIndex === -1) throw new Error('Game not found');
        if (this.games[gameIndex].user_id !== userId) throw new Error('Unauthorized');

        const game = this.games[gameIndex];
        this.games.splice(gameIndex, 1);

        this.guesses = this.guesses.filter(g => g.game_id !== gameId);
        this.save();

        return game;
    }

    // Guess endpoints
    async createGuess(gameId, guess) {
        await this.delay();

        const userId = this.getUserIdFromToken();
        if (!userId) throw new Error('Unauthorized');

        const game = this.games.find(g => g.id === gameId);
        if (!game) throw new Error('Game not found');
        if (game.user_id !== userId) throw new Error('Unauthorized');

        // Mock distance calculation (random for now)
        const distance = Math.floor(Math.random() * 100);

        const guessObj = {
            id: this.nextGuessId++,
            game_id: gameId,
            guess_word: guess,
            distance: distance,
            created_at: new Date().toISOString(),
        };

        this.guesses.push(guessObj);
        this.save();

        return { distance };
    }

    async getAllGuessByGame(gameId) {
        await this.delay();

        const userId = this.getUserIdFromToken();
        if (!userId) throw new Error('Unauthorized');

        const game = this.games.find(g => g.id === gameId);
        if (!game) throw new Error('Game not found');
        if (game.user_id !== userId) throw new Error('Unauthorized');

        return this.guesses
            .filter(g => g.game_id === gameId)
            .sort((a, b) => new Date(a.created_at) - new Date(b.created_at));
    }

    async getGuessById(gameId, guessId) {
        await this.delay();

        const userId = this.getUserIdFromToken();
        if (!userId) throw new Error('Unauthorized');

        const game = this.games.find(g => g.id === gameId);
        if (!game) throw new Error('Game not found');
        if (game.user_id !== userId) throw new Error('Unauthorized');

        const guess = this.guesses.find(g => g.id === guessId && g.game_id === gameId);
        if (!guess) throw new Error('Guess not found');

        return guess;
    }

    // Clear all mock data
    clearAllData() {
        this.users = [];
        this.games = [];
        this.guesses = [];
        this.nextGameId = 1;
        this.nextGuessId = 1;
        localStorage.removeItem('mock_users');
        localStorage.removeItem('mock_games');
        localStorage.removeItem('mock_guesses');
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        localStorage.removeItem('currentGameId');
    }
}

export const mockApi = new MockApiClient();
export default mockApi;