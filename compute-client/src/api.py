from fastapi import FastAPI, HTTPException, Request
from contextlib import asynccontextmanager
import sys
import asyncio
import logging
import time
import uuid

from model import WordSimilarityModel
from schemas import (
    StartGameRequest, StartGameResponse,
    GuessRequest, GuessResponse,
    HealthResponse
)
from monitor import ResourceMonitor
from logging_config import setup_logging

logger = logging.getLogger(__name__)

word_model = WordSimilarityModel()
monitor = None


async def periodic_monitoring(interval: int = 60):
    logger.info(f"Starting periodic monitoring (interval: {interval}s)")
    try:
        while True:
            monitor.log_metrics(context="Periodic")
            await asyncio.sleep(interval)
    except asyncio.CancelledError:
        logger.info("Periodic monitoring stopped")
        raise


@asynccontextmanager
async def lifespan(app: FastAPI):
    global monitor

    setup_logging(service_name="compute-client", log_dir="logs")
    logger.info("=" * 60)
    logger.info("Starting compute-client service")
    logger.info("=" * 60)

    monitor = ResourceMonitor()

    monitoring_task = asyncio.create_task(periodic_monitoring(interval=60))

    try:
        monitor.log_metrics(context="Startup")
        word_model.load_model("en", "glove-twitter-25") #glove-twitter-25  word2vec-google-news-300
        monitor.log_metrics(context="Model Loaded")

        # Warm up the model to speed up first query
        logger.info("Warming up model...")
        word_model.warm_up("en")
        monitor.log_metrics(context="Model Warmed Up")

        logger.info("Application started successfully")
    except Exception as e:
        logger.error(f"Failed to start application: {e}")
        monitoring_task.cancel()
        sys.exit(1)

    yield

    logger.info("=" * 60)
    logger.info("Shutting down compute-client service")
    logger.info("=" * 60)

    monitoring_task.cancel()
    try:
        await monitoring_task
    except asyncio.CancelledError:
        pass

    try:
        word_model.cleanup()
        monitor.log_metrics(context="Shutdown")
        logger.info("Cleanup completed successfully")
    except Exception as e:
        logger.error(f"Error during cleanup: {e}")
    finally:
        logger.info("Application stopped")


app = FastAPI(lifespan=lifespan)


@app.middleware("http")
async def log_request_resources(request: Request, call_next):
    request_id = str(uuid.uuid4())[:8]

    request.state.request_id = request_id

    start_time = time.time()

    try:
        response = await call_next(request)

        duration = time.time() - start_time
        metrics = monitor.get_metrics(include_cpu=False)  # Skip CPU for faster response

        logger.info(
            f"[{request_id}] {request.method} {request.url.path} | "
            f"Status: {response.status_code} | "
            f"Duration: {duration:.3f}s | "
            f"Memory: {metrics['memory_mb']:.1f}MB"
        )

        response.headers["X-Request-ID"] = request_id

        return response

    except Exception as e:
        duration = time.time() - start_time
        logger.error(
            f"[{request_id}] {request.method} {request.url.path} | "
            f"ERROR: {str(e)} | Duration: {duration:.3f}s"
        )
        raise


@app.post("/start-game", response_model=StartGameResponse)
async def start_game(request: StartGameRequest):
    language = request.language
    secret_word = request.secret_word.lower().strip()

    if not word_model.is_language_supported(language):
        raise HTTPException(
            status_code=503,
            detail=f"Language '{language}' not supported"
        )

    if not word_model.word_in_vocabulary(secret_word, language):
        raise HTTPException(
            status_code=400,
            detail="Secret word not in vocabulary"
        )

    try:
        rankings, calc_time = word_model.calculate_distance(
            secret_word, language, request.top_n
        )
    except Exception as e:
        raise HTTPException(
            status_code=500,
            detail=f"Failed to calculate: {str(e)}"
        )

    return StartGameResponse(
        calculation_time=calc_time,
        hint_word="my hint"
    )


@app.post("/guess", response_model=GuessResponse)
async def make_guess(request: GuessRequest):
    guess = request.guess.lower().strip()
    secret_word = request.secret_word

    try:
        distance, found = word_model.get_guess_distance(secret_word, guess)
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))

    return GuessResponse(distance=distance)


@app.get("/health", response_model=HealthResponse)
async def health_check():
    return HealthResponse(
        status="healthy",
        loaded_languages=list(word_model.models.keys())
    )
 

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=5000)