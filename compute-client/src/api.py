from fastapi import FastAPI, HTTPException
from contextlib import asynccontextmanager
import sys

from model import WordSimilarityModel
from schemas import (
    StartGameRequest, StartGameResponse,
    GuessRequest, GuessResponse,
    HealthResponse
)

word_model = WordSimilarityModel()


@asynccontextmanager
async def lifespan(app: FastAPI):
    print("Starting application...")
    try:
        word_model.load_model("en", "glove-twitter-25")
    except Exception as e:
        print(f"Failed to load model: {e}")
        sys.exit(1)

    yield

    print("Shutting down gracefully...")
    try:
        word_model.cleanup()
        print("Cleanup completed")
    except Exception as e:
        print(f"Error during cleanup: {e}")
    finally:
        print("Application stopped")


app = FastAPI(lifespan=lifespan)


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
        rankings, calc_time = word_model.calculate_rankings(
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
        rank, found = word_model.get_guess_rank(secret_word, guess)
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))

    return GuessResponse(rank=rank, found=found)


@app.get("/health", response_model=HealthResponse)
async def health_check():
    return HealthResponse(
        status="healthy",
        loaded_languages=list(word_model.models.keys())
    )


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=5000)