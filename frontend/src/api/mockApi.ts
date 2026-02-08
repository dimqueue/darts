// TODO: Consider replacing with MSW (Mock Service Worker) for more realistic mocking
// TODO: Move mockNames/mockUsernames arrays to a separate fixtures file to reduce this file's size (~462 lines)
import { getRandomWord, calculateDistance } from './mockData';
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
    type LeaderboardEntry,
    type RankData,
    type PublicProfileData,
    type ApiResponse,
} from '../types/api';
import type { Game, Guess, GuessResult } from '../types/game';
import { STORAGE_KEYS } from '../config/constants';

interface MockUser {
    id: number;
    username: string;
    password: string;
    name: string;
    created_at: string;
}

class MockApiClient {
    private users: Map<string, MockUser> = new Map();
    private games: Map<number, Game> = new Map();
    private guesses: Map<number, Guess[]> = new Map();
    private targetWords: Map<number, string> = new Map();
    private nextUserId = 1;
    private nextGameId = 1;
    private nextGuessId = 1;

    private delay(ms: number = 300): Promise<void> {
        return new Promise((resolve) => setTimeout(resolve, ms));
    }

    async signUp(username: string, password: string, name: string = ''): Promise<SignUpResponse> {
        await this.delay();
        if (this.users.has(username)) {
            throw new ApiError('Username already exists', 'USER_EXISTS', 409);
        }
        const user: MockUser = {
            id: this.nextUserId++,
            username,
            password,
            name: name || username,
            created_at: new Date().toISOString(),
        };
        this.users.set(username, user);
        return { message: 'User created successfully' };
    }

    async signIn(username: string, password: string): Promise<AuthResponse> {
        await this.delay();
        const user = this.users.get(username);
        if (!user || user.password !== password) {
            throw new ApiError('Invalid credentials', 'UNAUTHORIZED', 401);
        }
        return { token: `mock-token-${user.id}`, user_id: user.id };
    }

    async createGame(language: string): Promise<Game> {
        await this.delay();
        const game: Game = {
            id: this.nextGameId++,
            language,
            status: 'in_progress',
            started_at: new Date().toISOString(),
        };
        this.games.set(game.id, game);
        this.guesses.set(game.id, []);
        const targetWord = getRandomWord(language);
        this.targetWords.set(game.id, targetWord);
        return game;
    }

    async getAllGames(): Promise<ApiResponse<Game[]>> {
        await this.delay();
        return { data: Array.from(this.games.values()) };
    }

    async getGameById(gameId: number): Promise<Game> {
        await this.delay();
        const game = this.games.get(gameId);
        if (!game) throw new ApiError('Game not found', 'GAME_NOT_FOUND', 404);
        return game;
    }

    async updateGame(gameId: number, data: Partial<Game>): Promise<Game> {
        await this.delay();
        const game = this.games.get(gameId);
        if (!game) throw new ApiError('Game not found', 'GAME_NOT_FOUND', 404);
        Object.assign(game, data);
        return game;
    }

    async deleteGame(gameId: number): Promise<{ message: string }> {
        await this.delay();
        this.games.delete(gameId);
        this.guesses.delete(gameId);
        this.targetWords.delete(gameId);
        return { message: 'Game deleted' };
    }

    async abandonGame(gameId: number): Promise<{ message: string }> {
        await this.delay();
        const game = this.games.get(gameId);
        if (!game) throw new ApiError('Game not found', 'GAME_NOT_FOUND', 404);
        game.status = 'abandoned';
        return { message: 'Game abandoned' };
    }

    async createGuess(gameId: number, guess: string): Promise<GuessResult> {
        await this.delay();
        const gameGuesses = this.guesses.get(gameId) || [];
        const targetWord = this.targetWords.get(gameId) || 'ocean';
        const rank = calculateDistance(guess, targetWord);

        const found = guess.toLowerCase() === targetWord.toLowerCase();
        const inVocabulary = rank !== -1;

        const newGuess: Guess = {
            id: this.nextGuessId++,
            game_id: gameId,
            guess_word: guess,
            distance: Math.max(0, rank),
            created_at: new Date().toISOString(),
        };
        gameGuesses.push(newGuess);
        this.guesses.set(gameId, gameGuesses);

        return {
            rank: Math.max(0, rank),
            found: found,
            in_vocabulary: inVocabulary,
        };
    }

    async getAllGuessByGame(gameId: number): Promise<ApiResponse<Guess[]>> {
        await this.delay();
        return { data: this.guesses.get(gameId) || [] };
    }

    async getGuessById(gameId: number, guessId: number): Promise<Guess> {
        await this.delay();
        const gameGuesses = this.guesses.get(gameId) || [];
        const guess = gameGuesses.find((g) => g.id === guessId);
        if (!guess) throw new ApiError('Guess not found', 'NOT_FOUND', 404);
        return guess;
    }

    async getActiveGame(): Promise<ApiResponse<Game | null>> {
        await this.delay();
        const activeGame = Array.from(this.games.values()).find((g) => g.status === 'in_progress');
        return { data: activeGame || null };
    }

    async getMyProfile(): Promise<ProfileData> {
        await this.delay();

        let currentUser: { id?: number; username?: string; name?: string; email?: string } | null =
            null;
        try {
            const savedUser = localStorage.getItem(STORAGE_KEYS.USER);
            if (savedUser) {
                currentUser = JSON.parse(savedUser);
            }
        } catch {
            // Ignore JSON parse errors
        }

        return {
            id: currentUser?.id || 1,
            username: currentUser?.username || 'testuser',
            name: currentUser?.name || currentUser?.username || 'Test User',
            email: currentUser?.email || 'test@example.com',
            avatar_url: null,
            bio: 'Hello, I love playing Darts!',
            country_code: 'UA',
            total_games: 42,
            total_wins: 28,
            total_losses: 14,
            current_win_streak: 3,
            best_win_streak: 7,
            average_guesses: 4.5,
            total_score: 1250,
        };
    }

    async updateMyProfile(data: ProfileUpdateData): Promise<{ message: string }> {
        await this.delay();
        return { message: 'Profile updated', ...data };
    }

    async getMySettings(): Promise<SettingsData> {
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

    async updateMySettings(data: Partial<SettingsData>): Promise<{ message: string }> {
        await this.delay();
        return { message: 'Settings updated', ...data };
    }

    async getMyStatistics(): Promise<StatisticsData> {
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

    async getMyLanguageStats(): Promise<ApiResponse<LanguageStatsData[]>> {
        await this.delay();
        return {
            data: [
                {
                    language: 'en',
                    games_played: 30,
                    games_won: 20,
                    average_guesses: 4.2,
                    total_score: 900,
                },
                {
                    language: 'ua',
                    games_played: 12,
                    games_won: 8,
                    average_guesses: 5.1,
                    total_score: 350,
                },
            ],
        };
    }

    async getLeaderboard(
        type: string = 'global',
        _params: { limit?: number; offset?: number } = {},
        language: string | null = null
    ): Promise<LeaderboardResponse> {
        await this.delay();

        let currentUser: { id?: number; username?: string; name?: string } | null = null;
        try {
            const savedUser = localStorage.getItem(STORAGE_KEYS.USER);
            if (savedUser) {
                currentUser = JSON.parse(savedUser);
            }
        } catch {
            // Ignore
        }

        const mockNames = [
            'Alex Storm',
            'Maria Chen',
            'John Smith',
            'Emma Wilson',
            'Oleksandr Petrov',
            'Sophie Martin',
            'Michael Brown',
            'Anna Kovalenko',
            'David Lee',
            'Julia Schmidt',
            'James Taylor',
            'Katya Bondarenko',
            'Robert Garcia',
            'Nina Ivanova',
            'William Davis',
            'Elena Moroz',
            'Daniel Miller',
            'Oksana Shevchenko',
            'Christopher Moore',
            'Tetiana Lysenko',
            'Matthew Anderson',
            'Iryna Tkachenko',
            'Andrew Jackson',
            'Viktoria Melnyk',
            'Joshua White',
            'Natalia Kravchenko',
            'Ryan Harris',
            'Svitlana Boyko',
            'Brandon Clark',
            'Daryna Savchenko',
            'Kevin Lewis',
            'Alina Rudenko',
            'Jason Robinson',
            'Maryna Kozak',
            'Justin Walker',
            'Yulia Polishchuk',
            'Eric Hall',
            'Olena Marchenko',
            'Adam Young',
            'Karina Stepanenko',
            'Tyler King',
            'Larysa Hryhorenko',
            'Aaron Wright',
            'Veronika Fedorova',
            'Nicholas Scott',
            'Anastasia Zaitseva',
            'Patrick Green',
            'Diana Romanova',
            'Sean Baker',
            'Polina Sokolova',
        ];

        const mockUsernames = [
            'alexstorm',
            'mariachen',
            'johnsmith',
            'emmawilson',
            'oleksandrp',
            'sophiem',
            'mikebrown',
            'annakov',
            'davidlee',
            'juliaschmidt',
            'jamest',
            'katyab',
            'robertg',
            'ninaiv',
            'williamd',
            'elenam',
            'danmiller',
            'oksanas',
            'chrismoore',
            'tetianal',
            'mattanderson',
            'irynat',
            'andrewj',
            'viktoriamelnyk',
            'joshwhite',
            'nataliak',
            'ryanharris',
            'svitlanab',
            'brandonc',
            'darynas',
            'kevinlewis',
            'alinar',
            'jasonr',
            'marynakozak',
            'justinw',
            'yuliap',
            'erichall',
            'olenam',
            'adamyoung',
            'karinas',
            'tylerk',
            'larysah',
            'aaronwright',
            'veronikaf',
            'nicholasscott',
            'anastasiaz',
            'patrickg',
            'dianar',
            'seanbaker',
            'polinas',
        ];

        const ranksByType: Record<string, number> = {
            global: 23,
            daily: 2,
            weekly: 3,
            monthly: 18,
        };
        const currentUserRank = ranksByType[type] || 23;

        const mockUsers: LeaderboardEntry[] = Array.from({ length: 50 }, (_, i) => {
            const rank = i + 1;

            if (rank === currentUserRank && currentUser) {
                return {
                    rank: rank,
                    user_id: currentUser.id || 999,
                    username: currentUser.username || 'user',
                    name: currentUser.name || currentUser.username,
                    avatar_url: null,
                    country_code: 'UA',
                    total_score: 1500 - (rank - 1) * 25,
                    total_wins: 50 - (rank - 1),
                    total_games: 60 - (rank - 1),
                    best_win_streak: 10 - Math.floor((rank - 1) / 5),
                    average_guesses: 3.5 + (rank - 1) * 0.1,
                    win_rate: Number((((50 - (rank - 1)) / (60 - (rank - 1))) * 100).toFixed(1)),
                };
            }

            return {
                rank: rank,
                user_id: i + 100,
                username: mockUsernames[i] || `player${i + 1}`,
                name: mockNames[i] || `Player ${i + 1}`,
                avatar_url: null,
                country_code: ['US', 'GB', 'DE', 'FR', 'UA'][i % 5],
                total_score: 1500 - i * 25,
                total_wins: 50 - i,
                total_games: 60 - i,
                best_win_streak: 10 - Math.floor(i / 5),
                average_guesses: 3.5 + i * 0.1,
                win_rate: Number((((50 - i) / (60 - i)) * 100).toFixed(1)),
            };
        });

        return {
            leaderboard_type: type,
            language: language,
            users: mockUsers,
            total: 150,
            current_user_rank: currentUserRank,
        };
    }

    async getMyRank(): Promise<RankData> {
        await this.delay();
        return {
            global_rank: 23,
            daily_rank: 2,
            weekly_rank: 3,
            monthly_rank: 18,
        };
    }

    async getPublicProfile(username: string): Promise<PublicProfileData> {
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

    clearAllData(): void {
        this.users.clear();
        this.games.clear();
        this.guesses.clear();
        this.targetWords.clear();
        localStorage.removeItem(STORAGE_KEYS.TOKEN);
        localStorage.removeItem(STORAGE_KEYS.USER);
        localStorage.removeItem(STORAGE_KEYS.THEME);
        localStorage.removeItem(STORAGE_KEYS.DARK_MODE);
    }
}

export default MockApiClient;
