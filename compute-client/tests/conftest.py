import sys
from pathlib import Path
from unittest.mock import MagicMock, patch

import pytest

# Add src to path so tests can import from it
src_path = Path(__file__).parent.parent / "src"
sys.path.insert(0, str(src_path))


@pytest.fixture
def mock_word_model():
    #Mock WordSimilarityModel to avoid loading real ML models
    from models.word_similarity import WordSimilarityModel

    model = WordSimilarityModel(rankings_cache=None)

    vocab = {"cat", "dog", "apple", "secret", "hello", "world", "kitten", "pet", "animal", "feline"}

    class MockGensimModel:
        def __init__(self, vocabulary):
            self.vocab = vocabulary
            self.index_to_key = list(vocabulary)

        def __getitem__(self, word):
            if word in self.vocab:
                return [0.1] * 25  # fake vector
            raise KeyError(word)

        def most_similar(self, word, topn=10):
            return [
                ("dog", 0.9),
                ("kitten", 0.85),
                ("pet", 0.8),
                ("animal", 0.75),
                ("feline", 0.7),
            ][:topn]

    model.models["en"] = MockGensimModel(vocab)

    return model


@pytest.fixture
def mock_rankings_cache():
    #Mock RankingsCache with fake Redis
    try:
        import fakeredis

        fake_redis = fakeredis.FakeRedis(decode_responses=True)

        mock_redis_client = MagicMock()
        mock_redis_client.get_client.return_value = fake_redis
        mock_redis_client.is_connected.return_value = True

        from cache.rankings_cache import RankingsCache
        cache = RankingsCache(mock_redis_client, ttl=3600)

        return cache
    except ImportError:
        pytest.skip("fakeredis not installed")


@pytest.fixture
def sample_rankings():
    #Sample rankings data for testing
    return {
        "secret": 1,
        "mystery": 2,
        "hidden": 3,
        "unknown": 4,
        "puzzle": 5,
    }


@pytest.fixture
def mock_grpc_context():
    #Mock gRPC context for servicer tests.
    context = MagicMock()
    context.peer.return_value = "test-peer:12345"
    context.set_code = MagicMock()
    context.set_details = MagicMock()
    return context


@pytest.fixture
def mock_redis_client():
    #Mock Redis client
    try:
        import fakeredis

        fake_redis = fakeredis.FakeRedis(decode_responses=True)

        mock_client = MagicMock()
        mock_client.get_client.return_value = fake_redis
        mock_client.is_connected.return_value = True
        mock_client._client = fake_redis

        return mock_client
    except ImportError:
        pytest.skip("fakeredis not installed")
