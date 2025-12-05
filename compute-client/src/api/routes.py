"""API route handlers"""
import logging

from fastapi import APIRouter, HTTPException, Request

from config.schemas import (
    StartGameRequest, StartGameResponse,
    GuessRequest, GuessResponse,
    HealthResponse
)

logger = logging.getLogger(__name__)

router = APIRouter()


@router.post("/start-game", response_model=StartGameResponse)
async def start_game(request_data: StartGameRequest, request: Request):
    word_model = request.app.state.word_model
    language = request_data.language
    secret_word = request_data.secret_word.lower().strip()

    logger.info(
        f"[{getattr(request.state, 'request_id', 'N/A')}] Start game request: "
        f"language={language}, secret_word={secret_word}, top_n={request_data.top_n}"
    )

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
            secret_word, language, request_data.top_n
        )
        logger.info(
            f"[{getattr(request.state, 'request_id', 'N/A')}] Start game completed: "
            f"calculation_time={calc_time:.3f}s"
        )
    except Exception as e:
        logger.error(
            f"[{getattr(request.state, 'request_id', 'N/A')}] Start game failed: {str(e)}"
        )
        raise HTTPException(
            status_code=500,
            detail=f"Failed to calculate: {str(e)}"
        )

    return StartGameResponse(
        calculation_time=calc_time,
        hint_word="my hint"
    )


@router.post("/guess", response_model=GuessResponse)
async def make_guess(request_data: GuessRequest, request: Request):
    word_model = request.app.state.word_model
    guess = request_data.guess.lower().strip()
    secret_word = request_data.secret_word

    logger.info(
        f"[{getattr(request.state, 'request_id', 'N/A')}] Guess request: "
        f"guess={guess}, secret_word={secret_word}"
    )

    try:
        distance, found = word_model.get_guess_distance(secret_word, guess)
        logger.info(
            f"[{getattr(request.state, 'request_id', 'N/A')}] Guess result: "
            f"distance={distance}"
        )
    except ValueError as e:
        logger.error(
            f"[{getattr(request.state, 'request_id', 'N/A')}] Guess failed: {str(e)}"
        )
        raise HTTPException(status_code=404, detail=str(e))

    return GuessResponse(distance=distance)


@router.get("/health", response_model=HealthResponse)
async def health_check(request: Request):
    word_model = request.app.state.word_model
    return HealthResponse(
        status="healthy",
        loaded_languages=list(word_model.models.keys())
    )