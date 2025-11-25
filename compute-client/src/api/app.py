import asyncio
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI

from config.settings import Config
from models.word_similarity import WordSimilarityModel
from utils.monitor import ResourceMonitor
from .middleware import log_request_resources
from .routes import router

logger = logging.getLogger(__name__)


def create_lifespan(word_model: WordSimilarityModel, monitor: ResourceMonitor):
    @asynccontextmanager
    async def lifespan(app: FastAPI):
        app.state.word_model = word_model
        app.state.monitor = monitor

        logger.info("Starting periodic monitoring...")
        monitoring_task = asyncio.create_task(
            periodic_monitoring(monitor, Config.MONITORING_INTERVAL)
        )

        logger.info("HTTP server ready")

        yield

        logger.info("Shutting down HTTP server...")
        monitoring_task.cancel()
        try:
            await monitoring_task
        except asyncio.CancelledError:
            pass

        logger.info("HTTP server stopped")

    return lifespan


async def periodic_monitoring(monitor: ResourceMonitor, interval: int):
    logger.info(f"Periodic monitoring started (interval: {interval}s)")
    try:
        while True:
            monitor.log_metrics(context="Periodic")
            await asyncio.sleep(interval)
    except asyncio.CancelledError:
        logger.info("Periodic monitoring stopped")
        raise


def create_app(word_model: WordSimilarityModel, monitor: ResourceMonitor) -> FastAPI:
    app = FastAPI(
        title="Compute Client API",
        description="Word similarity computation service",
        lifespan=create_lifespan(word_model, monitor)
    )

    app.middleware("http")(log_request_resources)

    app.include_router(router)

    return app