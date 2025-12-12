"""Tests for WordSimilarityModel."""
import pytest
from unittest.mock import MagicMock, patch


class TestWordSimilarityModel:
    """Tests for the WordSimilarityModel class."""

    def test_is_language_supported_true(self, mock_word_model):
        """Test supported language returns True."""
        assert mock_word_model.is_language_supported("en") is True

    def test_is_language_supported_false(self, mock_word_model):
        """Test unsupported language returns False."""
        assert mock_word_model.is_language_supported("fr") is False
        assert mock_word_model.is_language_supported("") is False
        assert mock_word_model.is_language_supported("xyz") is False

    def test_word_in_vocabulary_exists(self, mock_word_model):
        """Test word exists in vocabulary."""
        # The mock model contains: cat, dog, apple, secret, hello, world
        assert mock_word_model.word_in_vocabulary("cat", "en") is True
        assert mock_word_model.word_in_vocabulary("dog", "en") is True
        assert mock_word_model.word_in_vocabulary("secret", "en") is True

    def test_word_in_vocabulary_not_exists(self, mock_word_model):
        """Test word not in vocabulary."""
        assert mock_word_model.word_in_vocabulary("xyznonexistent", "en") is False
        assert mock_word_model.word_in_vocabulary("notaword123", "en") is False

    def test_calculate_distance_returns_rankings(self, mock_word_model):
        """Test calculate_distance returns rankings dict and timing."""
        rankings, calc_time = mock_word_model.calculate_distance("cat", "en", 5)

        # Should have the secret word at rank 1
        assert rankings["cat"] == 1

        # Should have similar words
        assert "dog" in rankings
        assert rankings["dog"] == 2

        # calc_time should be a number (may be very small in tests)
        assert isinstance(calc_time, float)
        assert calc_time >= 0

    def test_calculate_distance_stores_in_memory(self, mock_word_model):
        """Test that rankings are cached in memory."""
        mock_word_model.calculate_distance("cat", "en", 5)

        # Should be stored in game_rankings
        assert "cat" in mock_word_model.game_rankings
        assert mock_word_model.game_rankings["cat"]["cat"] == 1

    def test_calculate_distance_uses_cache(self, mock_word_model, mock_rankings_cache):
        """Test that rankings are retrieved from cache if available."""
        # Set up the model with cache
        mock_word_model.rankings_cache = mock_rankings_cache

        # Pre-populate cache
        cached_rankings = {"secret": 1, "mystery": 2, "hidden": 3}
        mock_rankings_cache.set_rankings("en", "secret", cached_rankings)

        # Calculate should use cache
        rankings, calc_time = mock_word_model.calculate_distance("secret", "en", 10)

        # Should return cached data
        assert rankings == cached_rankings
        # Calc time should be 0 for cached results
        assert calc_time == 0.0

    def test_get_guess_distance_exact_match(self, mock_word_model):
        """Test guessing the exact secret word returns distance 1."""
        distance, found = mock_word_model.get_guess_distance("cat", "cat", "en")

        assert distance == 1
        assert found is True

    def test_get_guess_distance_not_in_vocabulary(self, mock_word_model):
        """Test guess not in vocabulary returns -1."""
        # First calculate rankings for the secret word
        mock_word_model.calculate_distance("cat", "en", 10)

        # Now check an unknown word
        distance, found = mock_word_model.get_guess_distance("cat", "xyznotaword", "en")

        assert distance == -1
        assert found is False

    def test_get_guess_distance_found_in_rankings(self, mock_word_model):
        """Test guess found in pre-calculated rankings."""
        # First calculate rankings for "cat"
        mock_word_model.calculate_distance("cat", "en", 10)

        # Now check distance for "dog" which should be in rankings
        distance, found = mock_word_model.get_guess_distance("cat", "dog", "en")

        assert distance == 2  # dog is at rank 2 (cat=1, dog=2, kitten=3, etc.)
        assert found is False

    def test_get_guess_distance_no_game_data_raises(self, mock_word_model):
        """Test that missing game data raises ValueError."""
        # Don't calculate rankings first
        mock_word_model.game_rankings = {}

        with pytest.raises(ValueError, match="No game data found"):
            mock_word_model.get_guess_distance("unknownsecret", "cat", "en")

    def test_get_guess_distance_not_in_top_n(self, mock_word_model):
        """Test word in vocabulary but not in top N returns 0."""
        # Calculate rankings with small topn
        mock_word_model.calculate_distance("cat", "en", 3)

        # "apple" is in vocabulary but not in the similar words list
        # The mock returns only: dog, kitten, pet, animal, feline
        distance, found = mock_word_model.get_guess_distance("cat", "apple", "en")

        # Should return 0 for word not in top N
        assert distance == 0
        assert found is False

    def test_cleanup_clears_models(self, mock_word_model):
        """Test cleanup removes all loaded models."""
        mock_word_model.game_rankings["test"] = {"test": 1}

        mock_word_model.cleanup()

        assert len(mock_word_model.models) == 0
        assert len(mock_word_model.game_rankings) == 0


class TestWordSimilarityModelWithCache:
    """Tests for WordSimilarityModel with RankingsCache integration."""

    def test_calculate_distance_caches_to_redis(self, mock_word_model, mock_rankings_cache):
        """Test that calculated rankings are cached to Redis."""
        mock_word_model.rankings_cache = mock_rankings_cache

        # Calculate rankings
        mock_word_model.calculate_distance("cat", "en", 5)

        # Verify it was cached
        assert mock_rankings_cache.exists("en", "cat") is True

        # Retrieve and verify
        cached = mock_rankings_cache.get_rankings("en", "cat")
        assert cached is not None
        assert cached["cat"] == 1

    def test_get_guess_distance_loads_from_redis(self, mock_word_model, mock_rankings_cache):
        """Test that rankings are loaded from Redis when not in memory."""
        mock_word_model.rankings_cache = mock_rankings_cache

        # Pre-populate Redis cache directly
        cached_rankings = {"secret": 1, "word1": 2, "word2": 3, "cat": 50}
        mock_rankings_cache.set_rankings("en", "secret", cached_rankings)

        # Clear in-memory cache
        mock_word_model.game_rankings = {}

        # This should load from Redis
        distance, found = mock_word_model.get_guess_distance("secret", "cat", "en")

        # Should find cat at rank 50
        assert distance == 50
        assert found is False

        # Should also be loaded into memory now
        assert "secret" in mock_word_model.game_rankings
