import gensim.downloader as api
from typing import Dict, Tuple
import time
import logging

from utils.progress import Spinner

logger = logging.getLogger(__name__)


class WordSimilarityModel:

    def __init__(self):
        self.models = {}
        self.game_rankings: Dict[str, Dict[str, int]] = {}

    def load_model(self, language: str, model_name: str = "glove-twitter-25"):
        logger.info(f"Loading model '{model_name}' for language: {language}")

        spinner = Spinner(f"Loading {model_name}")

        spinner.update()

        self.models[language] = api.load(model_name)

        spinner.finish(f"Model '{model_name}' loaded successfully")
        logger.info(f"Model loaded for {language}")

    def is_language_supported(self, language: str) -> bool:
        return language in self.models

    def word_in_vocabulary(self, word: str, language: str) -> bool:
        try:
            self.models[language][word]
            return True
        except KeyError:
            return False

    def calculate_distance(self, secret_word: str, language: str, top_n: int) -> Tuple[Dict[str, int], float]:

        start_time = time.time()

        model = self.models[language]
        similar_words = model.most_similar(secret_word, topn=top_n)

        rankings = {secret_word: 1}
        for rank, (word, similarity) in enumerate(similar_words, start=2):
            rankings[word.lower()] = rank

        calculation_time = time.time() - start_time

        self.game_rankings[secret_word] = rankings

        return rankings, calculation_time

    def get_guess_distance(self, secret_word: str, guess: str) -> Tuple[int, bool]:

        game_data = self.game_rankings.get(secret_word)

        if not game_data:
            raise ValueError(f"No game data found for secret word: {secret_word}")

        if guess == secret_word:
            return 1, True

        rank = game_data.get(guess, 0)
        return rank, False

    def warm_up(self, language: str):
        if language not in self.models:
            raise ValueError(f"Language '{language}' not loaded")

        model = self.models[language]

        vocab_words = list(model.index_to_key)
        if vocab_words:
            warm_up_word = vocab_words[0]
            logger.info(f"Warming up model with word: '{warm_up_word}'")

            spinner = Spinner("Warming up model")
            spinner.update()

            start_time = time.time()
            _ = model.most_similar(warm_up_word, topn=10)
            warm_up_time = time.time() - start_time

            spinner.finish(f"Warm-up completed in {warm_up_time:.3f}s")
        else:
            logger.warning("No words found in vocabulary for warm-up")

    def cleanup(self):
        for lang in list(self.models.keys()):
            del self.models[lang]
        self.models.clear()
        self.game_rankings.clear()