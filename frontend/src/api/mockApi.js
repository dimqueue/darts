class MockApiClient {
    constructor() {
        console.log('API Mode: MOCK (mock data)');
        this.users = new Map();
        this.games = new Map();
        this.guesses = new Map();
        this.nextUserId = 1;
        this.nextGameId = 1;
        this.nextGuessId = 1;
    }

    delay(ms = 300) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }


    async signUp(username, password, name = '') {
        await this.delay();
        if (this.users.has(username)) {
            throw new Error('Username already exists');
        }
        const user = {
            id: this.nextUserId++,
            username,
            password,
            name: name || username,
            created_at: new Date().toISOString(),
        };
        this.users.set(username, user);
        return { message: 'User created successfully' };
    }

    async signIn(username, password) {
        await this.delay();
        const user = this.users.get(username);
        if (!user || user.password !== password) {
            throw new Error('Invalid credentials');
        }
        return { token: `mock-token-${user.id}`, user_id: user.id };
    }

    async createGame(language) {
        await this.delay();
        const game = {
            id: this.nextGameId++,
            language,
            status: 'in_progress',
            started_at: new Date().toISOString(),
        };
        this.games.set(game.id, game);
        this.guesses.set(game.id, []);
        return game;
    }

    async getAllGames() {
        await this.delay();
        return { data: Array.from(this.games.values()) };
    }

    async getGameById(gameId) {
        await this.delay();
        const game = this.games.get(gameId);
        if (!game) throw new Error('Game not found');
        return game;
    }

    async updateGame(gameId, data) {
        await this.delay();
        const game = this.games.get(gameId);
        if (!game) throw new Error('Game not found');
        Object.assign(game, data);
        return game;
    }

    async deleteGame(gameId) {
        await this.delay();
        this.games.delete(gameId);
        this.guesses.delete(gameId);
        return { message: 'Game deleted' };
    }


    async createGuess(gameId, guess) {
        await this.delay();
        const gameGuesses = this.guesses.get(gameId) || [];
        const newGuess = {
            id: this.nextGuessId++,
            game_id: gameId,
            guess_word: guess,
            distance: Math.floor(Math.random() * 1000),
            created_at: new Date().toISOString(),
        };
        gameGuesses.push(newGuess);
        this.guesses.set(gameId, gameGuesses);
        return newGuess;
    }

    async getAllGuessByGame(gameId) {
        await this.delay();
        return { data: this.guesses.get(gameId) || [] };
    }

    async getGuessById(gameId, guessId) {
        await this.delay();
        const gameGuesses = this.guesses.get(gameId) || [];
        const guess = gameGuesses.find(g => g.id === guessId);
        if (!guess) throw new Error('Guess not found');
        return guess;
    }


    async getMyProfile() {
        await this.delay();
        return {
            id: 1,
            username: 'testuser',
            name: 'Test User',
            email: 'test@example.com',
            avatar_url: null,
            bio: 'Hello, I love playing Darts!',
            country_code: 'US',
            total_games: 42,
            total_wins: 28,
            total_losses: 14,
            current_win_streak: 3,
            best_win_streak: 7,
            average_guesses: 4.5,
            total_score: 1250,
        };
    }

    async updateMyProfile(data) {
        await this.delay();
        return { message: 'Profile updated', ...data };
    }

    async getMySettings() {
        await this.delay();
        return {
            preferred_language: 'en',
            theme: 'purple',
            sound_enabled: true,
            email_notifications: true,
            show_profile_public: true,
            show_stats_public: true,
        };
    }

    async updateMySettings(data) {
        await this.delay();
        return { message: 'Settings updated', ...data };
    }

    async getMyStatistics() {
        await this.delay();
        return {
            total_games: 42,
            total_wins: 28,
            total_losses: 14,
            current_win_streak: 3,
            best_win_streak: 7,
            total_guesses: 189,
            average_guesses: 4.5,
            fastest_win_seconds: 45,
            fewest_guesses_win: 2,
            total_score: 1250,
        };
    }

    async getMyLanguageStats() {
        await this.delay();
        return {
            data: [
                { language: 'en', games_played: 30, games_won: 20, average_guesses: 4.2, total_score: 900 },
                { language: 'ua', games_played: 12, games_won: 8, average_guesses: 5.1, total_score: 350 },
            ],
        };
    }


    async getLeaderboard(type = 'global', params = {}, language = null) {
        await this.delay();
        const mockUsers = Array.from({ length: 50 }, (_, i) => ({
            rank: i + 1,
            user_id: i + 1,
            username: `player${i + 1}`,
            name: `Player ${i + 1}`,
            avatar_url: null,
            country_code: ['US', 'GB', 'DE', 'FR', 'UA'][i % 5],
            total_score: 1500 - i * 25,
            total_wins: 50 - i,
            total_games: 60 - i,
            best_win_streak: 10 - Math.floor(i / 5),
            average_guesses: 3.5 + (i * 0.1),
            win_rate: ((50 - i) / (60 - i) * 100).toFixed(1),
        }));
        return {
            leaderboard_type: type,
            language: language,
            users: mockUsers,
            total: 150,
            current_user_rank: 23,
        };
    }

    async getMyRank() {
        await this.delay();
        return {
            global_rank: 23,
            daily_rank: 12,
            weekly_rank: 15,
            monthly_rank: 18,
        };
    }


    async getPublicProfile(username) {
        await this.delay();
        return {
            username,
            name: 'Public User',
            avatar_url: null,
            country_code: 'US',
            total_games: 30,
            total_wins: 20,
            win_rate: 66.7,
            best_win_streak: 5,
        };
    }

    clearAllData() {
        this.users.clear();
        this.games.clear();
        this.guesses.clear();
        localStorage.removeItem('token');
        localStorage.removeItem('user');
    }
}

export default MockApiClient;