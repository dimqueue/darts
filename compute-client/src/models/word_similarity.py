import gensim.downloader as api
from typing import Dict, Optional, Tuple, TYPE_CHECKING
import time
import logging

from utils.progress import Spinner

if TYPE_CHECKING:
    from cache import RankingsCache

logger = logging.getLogger(__name__)


class WordSimilarityModel:

    def __init__(self, rankings_cache: Optional["RankingsCache"] = None):
        self.models = {}
        self.game_rankings: Dict[str, Dict[str, int]] = {}
        self.rankings_cache = rankings_cache

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
        if self.rankings_cache:
            cached = self.rankings_cache.get_rankings(language, secret_word)
            if cached:
                logger.debug(f"Using cached rankings for '{secret_word}'")
                self.game_rankings[secret_word] = cached  # Also store in memory for fast access
                return cached, 0.0

        start_time = time.time()

        model = self.models[language]
        similar_words = model.most_similar(secret_word, topn=top_n)

        rankings = {secret_word: 1}
        for rank, (word, similarity) in enumerate(similar_words, start=2):
            rankings[word.lower()] = rank

        calculation_time = time.time() - start_time

        if self.rankings_cache:
            self.rankings_cache.set_rankings(language, secret_word, rankings)

        self.game_rankings[secret_word] = rankings

        return rankings, calculation_time

    def get_guess_distance(self, secret_word: str, guess: str, language: str = "en") -> Tuple[int, bool]:
        if guess == secret_word:
            return 1, True

        if not self.word_in_vocabulary(guess, language):
            logger.debug(f"Word '{guess}' not in vocabulary")
            return -1, False

        # Check in-memory cache first
        game_data = self.game_rankings.get(secret_word)

        # If not in memory, try Redis
        if not game_data and self.rankings_cache:
            game_data = self.rankings_cache.get_rankings(language, secret_word)
            if game_data:
                logger.debug(f"Loaded rankings from Redis for '{secret_word}'")
                self.game_rankings[secret_word] = game_data

        if not game_data:
            raise ValueError(f"No game data found for secret word: {secret_word}")

        rank = game_data.get(guess, 0)  # 0 = exists but not in top N
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