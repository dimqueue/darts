import json
import logging
from typing import Dict, Optional

from .redis_client import RedisClient

logger = logging.getLogger(__name__)

DEFAULT_TTL = 86400  # 24 hours


class RankingsCache:
    def __init__(self, redis_client: RedisClient, ttl: int = DEFAULT_TTL):
        self.redis_client = redis_client
        self.ttl = ttl

    def _make_key(self, language: str, secret_word: str) -> str:
        return f"game:{language}:{secret_word}"

    def get_rankings(self, language: str, secret_word: str) -> Optional[Dict[str, int]]:
        client = self.redis_client.get_client()
        if client is None:
            logger.warning("Redis not available, cannot get rankings from cache")
            return None

        key = self._make_key(language, secret_word)
        try:
            data = client.get(key)
            if data:
                logger.debug(f"Cache hit for {key}")
                return json.loads(data)
            logger.debug(f"Cache miss for {key}")
            return None
        except Exception as e:
            logger.error(f"Failed to get rankings from cache: {e}")
            return None

    def set_rankings(self, language: str, secret_word: str, rankings: Dict[str, int]) -> bool:
        client = self.redis_client.get_client()
        if client is None:
            logger.warning("Redis not available, cannot cache rankings")
            return False

        key = self._make_key(language, secret_word)
        try:
            client.setex(key, self.ttl, json.dumps(rankings))
            logger.debug(f"Cached rankings for {key} (TTL: {self.ttl}s)")
            return True
        except Exception as e:
            logger.error(f"Failed to cache rankings: {e}")
            return False

    def delete_rankings(self, language: str, secret_word: str) -> bool:
        client = self.redis_client.get_client()
        if client is None:
            return False

        key = self._make_key(language, secret_word)
        try:
            client.delete(key)
            logger.debug(f"Deleted cache for {key}")
            return True
        except Exception as e:
            logger.error(f"Failed to delete cache: {e}")
            return False

    def exists(self, language: str, secret_word: str) -> bool:
        client = self.redis_client.get_client()
        if client is None:
            return False

        key = self._make_key(language, secret_word)
        try:
            return client.exists(key) > 0
        except Exception as e:
            logger.error(f"Failed to check cache existence: {e}")
            return False