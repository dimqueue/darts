"""
Configuration management using environment variables.
"""
import os
from pathlib import Path
from typing import Literal

try:
    from dotenv import load_dotenv
    # Load root .env first (for shared config like Redis)
    root_env = Path(__file__).parent.parent.parent.parent / ".env"
    if root_env.exists():
        load_dotenv(root_env)
    # Load local .env (overrides root values if present)
    load_dotenv()
except ImportError:
    pass


class Config:
    SERVER_MODE: Literal["http", "grpc"] = os.getenv("SERVER_MODE", "http")

    HTTP_PORT: int = int(os.getenv("HTTP_PORT", "5000"))
    GRPC_PORT: int = int(os.getenv("GRPC_PORT", "50051"))
    HTTP_HOST: str = os.getenv("HTTP_HOST", "0.0.0.0")

    MODEL_NAME: str = os.getenv("MODEL_NAME", "glove-twitter-25")
    DEFAULT_LANGUAGE: str = os.getenv("DEFAULT_LANGUAGE", "en")

    MONITORING_INTERVAL: int = int(os.getenv("MONITORING_INTERVAL", "60"))

    LOG_DIR: str = os.getenv("LOG_DIR", "logs")
    LOG_LEVEL: str = os.getenv("LOG_LEVEL", "DEBUG").upper()

    # Redis configuration
    REDIS_ENABLED: bool = os.getenv("REDIS_ENABLED", "true").lower() == "true"
    REDIS_HOST: str = os.getenv("REDIS_HOST", "localhost")
    REDIS_PORT: int = int(os.getenv("REDIS_PORT", "6379"))
    REDIS_DB: int = int(os.getenv("REDIS_DB", "0"))
    REDIS_PASSWORD: str = os.getenv("REDIS_PASSWORD", "")
    REDIS_CACHE_TTL: int = int(os.getenv("REDIS_CACHE_TTL", "86400"))  # 24 hours

    @classmethod
    def validate(cls):
        if cls.SERVER_MODE not in ["http", "grpc"]:
            raise ValueError(f"Invalid SERVER_MODE: {cls.SERVER_MODE}. Must be 'http' or 'grpc'")

        if cls.HTTP_PORT < 1 or cls.HTTP_PORT > 65535:
            raise ValueError(f"Invalid HTTP_PORT: {cls.HTTP_PORT}")

        if cls.GRPC_PORT < 1 or cls.GRPC_PORT > 65535:
            raise ValueError(f"Invalid GRPC_PORT: {cls.GRPC_PORT}")

        if cls.MONITORING_INTERVAL < 1:
            raise ValueError(f"Invalid MONITORING_INTERVAL: {cls.MONITORING_INTERVAL}")

    @classmethod
    def display(cls):
        return {
            "SERVER_MODE": cls.SERVER_MODE,
            "HTTP_PORT": cls.HTTP_PORT,
            "GRPC_PORT": cls.GRPC_PORT,
            "HTTP_HOST": cls.HTTP_HOST,
            "MODEL_NAME": cls.MODEL_NAME,
            "DEFAULT_LANGUAGE": cls.DEFAULT_LANGUAGE,
            "MONITORING_INTERVAL": cls.MONITORING_INTERVAL,
            "LOG_DIR": cls.LOG_DIR,
            "REDIS_ENABLED": cls.REDIS_ENABLED,
            "REDIS_HOST": cls.REDIS_HOST,
            "REDIS_PORT": cls.REDIS_PORT,
            "REDIS_DB": cls.REDIS_DB,
            "REDIS_CACHE_TTL": cls.REDIS_CACHE_TTL,
        }


Config.validate()
