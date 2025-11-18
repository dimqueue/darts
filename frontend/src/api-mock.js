import { mockApi } from './mockApi.js';

class MockApiWrapper {
    constructor() {
        console.log('API Mode: MOCK (mock data)');
        this.mockApi = mockApi;
    }

    async signUp(username, password, name) {
        return this.mockApi.signUp(username, password, name);
    }

    async signIn(username, password) {
        return this.mockApi.signIn(username, password);
    }

    async createGame(language) {
        return this.mockApi.createGame(language);
    }

    async getAllGames() {
        return this.mockApi.getAllGames();
    }

    async getGameById(gameId) {
        return this.mockApi.getGameById(gameId);
    }

    async updateGame(gameId, data) {
        return this.mockApi.updateGame(gameId, data);
    }

    async deleteGame(gameId) {
        return this.mockApi.deleteGame(gameId);
    }

    async createGuess(gameId, guess) {
        return this.mockApi.createGuess(gameId, guess);
    }

    async getAllGuessByGame(gameId) {
        return this.mockApi.getAllGuessByGame(gameId);
    }

    async getGuessById(gameId, guessId) {
        return this.mockApi.getGuessById(gameId, guessId);
    }

    async clearAllData() {
        this.mockApi.clearAllData();
    }
}

export default MockApiWrapper;