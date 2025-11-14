import gensim.downloader as api
from typing import Dict, Tuple
import time

class WordSimilarityModel:

    def __init__(self):
        self.models = {}
        self.game_rankings: Dict[str, Dict[str, int]] = {}

    def load_model(self, language: str, model_name: str = "glove-twitter-25"):
        print(f"Loading model for language: {language}")
        self.models[language] = api.load(model_name)
        print(f"Model loaded successfully for {language}")

    def is_language_supported(self, language: str) -> bool:
        return language in self.models

    def word_in_vocabulary(self, word: str, language: str) -> bool:
        try:
            self.models[language][word]
            return True
        except KeyError:
            return False

    def calculate_rankings(self, secret_word: str, language: str, top_n: int) -> Tuple[Dict[str, int], float]:

        start_time = time.time()

        model = self.models[language]
        similar_words = model.most_similar(secret_word, topn=top_n)

        rankings = {secret_word: 1}
        for rank, (word, similarity) in enumerate(similar_words, start=2):
            rankings[word.lower()] = rank

        calculation_time = time.time() - start_time

        self.game_rankings[secret_word] = rankings

        return rankings, calculation_time

    def get_guess_rank(self, secret_word: str, guess: str) -> Tuple[int, bool]:

        game_data = self.game_rankings.get(secret_word)

        if not game_data:
            raise ValueError(f"No game data found for secret word: {secret_word}")

        if guess == secret_word:
            return 1, True

        rank = game_data.get(guess, 0)
        return rank, False

    def cleanup(self):
        for lang in list(self.models.keys()):
            del self.models[lang]
        self.models.clear()
        self.game_rankings.clear()