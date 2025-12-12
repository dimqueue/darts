"""Tests for RankingsCache."""
import json
import pytest
from unittest.mock import MagicMock


class TestRankingsCache:
    """Tests for the RankingsCache class."""

    @pytest.fixture
    def rankings_cache(self, mock_redis_client):
        """Create a RankingsCache with mock Redis."""
        from cache.rankings_cache import RankingsCache
        return RankingsCache(mock_redis_client, ttl=3600)

    def test_make_key(self, rankings_cache):
        """Test cache key generation."""
        key = rankings_cache._make_key("en", "secret")
        assert key == "game:en:secret"

        key = rankings_cache._make_key("ua", "слово")
        assert key == "game:ua:слово"

    def test_set_and_get_rankings(self, rankings_cache, sample_rankings):
        """Test setting and retrieving cached rankings."""
        # Set rankings
        result = rankings_cache.set_rankings("en", "secret", sample_rankings)
        assert result is True

        # Get rankings
        retrieved = rankings_cache.get_rankings("en", "secret")
        assert retrieved is not None
        assert retrieved == sample_rankings

    def test_get_nonexistent_key(self, rankings_cache):
        """Test getting non-existent key returns None."""
        result = rankings_cache.get_rankings("en", "nonexistent")
        assert result is None

    def test_exists_true(self, rankings_cache, sample_rankings):
        """Test exists returns True for cached key."""
        rankings_cache.set_rankings("en", "secret", sample_rankings)

        assert rankings_cache.exists("en", "secret") is True

    def test_exists_false(self, rankings_cache):
        """Test exists returns False for missing key."""
        assert rankings_cache.exists("en", "nonexistent") is False

    def test_delete_rankings(self, rankings_cache, sample_rankings):
        """Test deleting cached entry."""
        # Set rankings
        rankings_cache.set_rankings("en", "secret", sample_rankings)
        assert rankings_cache.exists("en", "secret") is True

        # Delete
        result = rankings_cache.delete_rankings("en", "secret")
        assert result is True

        # Verify deleted
        assert rankings_cache.exists("en", "secret") is False
        assert rankings_cache.get_rankings("en", "secret") is None

    def test_delete_nonexistent_key(self, rankings_cache):
        """Test deleting non-existent key doesn't error."""
        result = rankings_cache.delete_rankings("en", "nonexistent")
        assert result is True

    def test_rankings_with_large_data(self, rankings_cache):
        """Test caching large rankings data (10000 words)."""
        large_rankings = {f"word{i}": i for i in range(1, 10001)}
        large_rankings["secret"] = 1

        result = rankings_cache.set_rankings("en", "secret", large_rankings)
        assert result is True

        retrieved = rankings_cache.get_rankings("en", "secret")
        assert retrieved is not None
        assert len(retrieved) == 10001
        assert retrieved["secret"] == 1
        assert retrieved["word10000"] == 10000

    def test_rankings_json_serialization(self, rankings_cache):
        """Test that rankings are properly JSON serialized/deserialized."""
        rankings = {
            "word1": 1,
            "word2": 2,
            "special-word": 100,
            "word_with_underscore": 500,
        }

        rankings_cache.set_rankings("en", "test", rankings)
        retrieved = rankings_cache.get_rankings("en", "test")

        assert retrieved == rankings
        assert isinstance(retrieved["word1"], int)


class TestRankingsCacheNoRedis:
    """Tests for RankingsCache behavior when Redis is unavailable."""

    @pytest.fixture
    def cache_no_redis(self):
        """Create a RankingsCache with unavailable Redis."""
        from cache.rankings_cache import RankingsCache

        mock_client = MagicMock()
        mock_client.get_client.return_value = None
        mock_client.is_connected.return_value = False

        return RankingsCache(mock_client, ttl=3600)

    def test_get_rankings_no_redis(self, cache_no_redis):
        """Test get_rankings returns None when Redis unavailable."""
        result = cache_no_redis.get_rankings("en", "secret")
        assert result is None

    def test_set_rankings_no_redis(self, cache_no_redis):
        """Test set_rankings returns False when Redis unavailable."""
        result = cache_no_redis.set_rankings("en", "secret", {"word": 1})
        assert result is False

    def test_exists_no_redis(self, cache_no_redis):
        """Test exists returns False when Redis unavailable."""
        result = cache_no_redis.exists("en", "secret")
        assert result is False

    def test_delete_rankings_no_redis(self, cache_no_redis):
        """Test delete_rankings returns False when Redis unavailable."""
        result = cache_no_redis.delete_rankings("en", "secret")
        assert result is False


class TestRankingsCacheErrors:
    """Tests for RankingsCache error handling."""

    @pytest.fixture
    def cache_with_errors(self):
        """Create a RankingsCache with error-throwing Redis."""
        from cache.rankings_cache import RankingsCache

        mock_redis = MagicMock()
        mock_redis.get.side_effect = Exception("Redis error")
        mock_redis.setex.side_effect = Exception("Redis error")
        mock_redis.delete.side_effect = Exception("Redis error")
        mock_redis.exists.side_effect = Exception("Redis error")

        mock_client = MagicMock()
        mock_client.get_client.return_value = mock_redis

        return RankingsCache(mock_client, ttl=3600)

    def test_get_rankings_handles_error(self, cache_with_errors):
        """Test get_rankings handles Redis errors gracefully."""
        result = cache_with_errors.get_rankings("en", "secret")
        assert result is None

    def test_set_rankings_handles_error(self, cache_with_errors):
        """Test set_rankings handles Redis errors gracefully."""
        result = cache_with_errors.set_rankings("en", "secret", {"word": 1})
        assert result is False

    def test_exists_handles_error(self, cache_with_errors):
        """Test exists handles Redis errors gracefully."""
        result = cache_with_errors.exists("en", "secret")
        assert result is False

    def test_delete_handles_error(self, cache_with_errors):
        """Test delete_rankings handles Redis errors gracefully."""
        result = cache_with_errors.delete_rankings("en", "secret")
        assert result is False
