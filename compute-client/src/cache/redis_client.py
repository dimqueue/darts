import logging
from typing import Optional
import redis

logger = logging.getLogger(__name__)


class RedisClient:
    _instance: Optional["RedisClient"] = None
    _client: Optional[redis.Redis] = None

    def __init__(self, host: str, port: int, db: int = 0, password: Optional[str] = None):
        self.host = host
        self.port = port
        self.db = db
        self.password = password
        self._client = None

    @classmethod
    def get_instance(cls) -> Optional["RedisClient"]:
        return cls._instance

    @classmethod
    def initialize(cls, host: str, port: int, db: int = 0, password: Optional[str] = None) -> "RedisClient":
        if cls._instance is None:
            cls._instance = cls(host, port, db, password)
            cls._instance.connect()
        return cls._instance

    def connect(self) -> bool:
        try:
            self._client = redis.Redis(
                host=self.host,
                port=self.port,
                db=self.db,
                password=self.password,
                decode_responses=True,
                socket_timeout=5,
                socket_connect_timeout=5,
                retry_on_timeout=True,
            )
            self._client.ping()
            logger.info(f"Connected to Redis at {self.host}:{self.port}")
            return True
        except redis.ConnectionError as e:
            logger.error(f"Failed to connect to Redis: {e}")
            self._client = None
            return False

    def is_connected(self) -> bool:
        if self._client is None:
            return False
        try:
            self._client.ping()
            return True
        except redis.ConnectionError:
            return False

    def get_client(self) -> Optional[redis.Redis]:
        return self._client

    def close(self):
        if self._client:
            self._client.close()
            self._client = None
            logger.info("Redis connection closed")