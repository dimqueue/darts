"""gRPC service implementation (servicer)"""
import grpc
import logging

from proto.compute.v1 import service_pb2
from proto.compute.v1 import service_pb2_grpc
from models.word_similarity import WordSimilarityModel

logger = logging.getLogger(__name__)


class ComputeServiceServicer(service_pb2_grpc.ComputeServiceServicer):

    def __init__(self, word_model: WordSimilarityModel):
        self.word_model = word_model
        logger.info("ComputeServiceServicer initialized")

    def StartGame(self, request, context):
        # Extract client info from context
        peer = context.peer()

        logger.info(
            f"StartGame request from {peer}: "
            f"language={request.language}, secret_word={request.secret_word}, "
            f"top_n={request.top_n if request.top_n > 0 else 10000}"
        )

        language = request.language
        secret_word = request.secret_word.lower().strip()
        top_n = request.top_n if request.top_n > 0 else 10000

        if not self.word_model.is_language_supported(language):
            logger.warning(f"Language '{language}' not supported (client: {peer})")
            context.set_code(grpc.StatusCode.FAILED_PRECONDITION)
            context.set_details(f"Language '{language}' not supported")
            return service_pb2.StartGameResponse()

        if not self.word_model.word_in_vocabulary(secret_word, language):
            logger.warning(f"Secret word '{secret_word}' not in vocabulary (client: {peer})")
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("Secret word not in vocabulary")
            return service_pb2.StartGameResponse()

        try:
            rankings, calc_time = self.word_model.calculate_distance(
                secret_word, language, top_n
            )

            response = service_pb2.StartGameResponse(
                calculation_time=calc_time,
                hint_word="my hint"
            )

            logger.info(
                f"StartGame completed: calculation_time={calc_time:.3f}s, "
                f"client={peer}"
            )
            return response

        except Exception as e:
            logger.error(f"StartGame failed: {e} (client: {peer})")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Failed to calculate: {str(e)}")
            return service_pb2.StartGameResponse()

    def MakeGuess(self, request, context):
        # Extract client info from context
        peer = context.peer()

        logger.info(
            f"MakeGuess request from {peer}: "
            f"guess={request.guess}, secret_word={request.secret_word}"
        )

        guess = request.guess.lower().strip()
        secret_word = request.secret_word
        language = request.language 

        try:
            distance, found = self.word_model.get_guess_distance(secret_word, guess, language)

            response = service_pb2.MakeGuessResponse(
                distance=distance
            )

            logger.info(f"MakeGuess result: distance={distance}, client={peer}")
            return response

        except ValueError as e:
            logger.error(f"MakeGuess failed: {e} (client: {peer})")
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(str(e))
            return service_pb2.MakeGuessResponse()

    def HealthCheck(self, request, context):
        # Extract client info from context
        peer = context.peer()

        logger.info(f"HealthCheck request from {peer}")

        response = service_pb2.HealthCheckResponse(
            status="healthy",
            loaded_languages=list(self.word_model.models.keys())
        )

        logger.info(f"HealthCheck response: status=healthy, languages={len(self.word_model.models)}, client={peer}")

        return response