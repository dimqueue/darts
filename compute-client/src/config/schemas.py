from pydantic import BaseModel
from typing import Optional

class StartGameRequest(BaseModel):
    language: str = "en"
    secret_word: str
    top_n: int = 10000

class StartGameResponse(BaseModel):
    calculation_time: float
    hint_word: Optional[str] = None

class GuessRequest(BaseModel):
    secret_word: str
    guess: str
    language: str = "en"

class GuessResponse(BaseModel):
    rank: int
    found: bool
    in_vocabulary: bool

class HealthResponse(BaseModel):
    status: str
    loaded_languages: list[str]