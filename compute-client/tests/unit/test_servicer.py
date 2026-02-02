"""Tests for gRPC ComputeServiceServicer."""
import pytest
from unittest.mock import MagicMock
import grpc


class TestComputeServiceServicer:
    """Tests for the ComputeServiceServicer gRPC handler."""

    @pytest.fixture
    def servicer(self, mock_word_model):
        """Create a servicer with mock word model."""
        from grpc_server.servicer import ComputeServiceServicer
        return ComputeServiceServicer(mock_word_model)

    @pytest.fixture
    def start_game_request(self):
        """Create a mock StartGame request."""
        from proto.compute.v1 import service_pb2
        return service_pb2.StartGameRequest(
            language="en",
            secret_word="cat",
            top_n=10000
        )

    @pytest.fixture
    def make_guess_request(self):
        """Create a mock MakeGuess request."""
        from proto.compute.v1 import service_pb2
        return service_pb2.MakeGuessRequest(
            guess="dog",
            secret_word="cat",
            language="en"
        )

    # --- StartGame Tests ---

    def test_start_game_success(self, servicer, start_game_request, mock_grpc_context):
        """Test successful game initialization."""
        response = servicer.StartGame(start_game_request, mock_grpc_context)

        # Should return a response with calculation time
        assert response.calculation_time >= 0
        assert response.hint_word == "my hint"

        # Context should not have error set
        mock_grpc_context.set_code.assert_not_called()

    def test_start_game_unsupported_language(self, servicer, mock_grpc_context):
        """Test error when language not supported."""
        from proto.compute.v1 import service_pb2

        request = service_pb2.StartGameRequest(
            language="fr",  # Not supported
            secret_word="bonjour",
            top_n=10000
        )

        response = servicer.StartGame(request, mock_grpc_context)

        # Should set error code
        mock_grpc_context.set_code.assert_called_once_with(grpc.StatusCode.FAILED_PRECONDITION)
        mock_grpc_context.set_details.assert_called_once()

    def test_start_game_word_not_in_vocab(self, servicer, mock_grpc_context):
        """Test error when target word not in vocabulary."""
        from proto.compute.v1 import service_pb2

        # Use a word that's definitely not in our mock vocabulary
        request = service_pb2.StartGameRequest(
            language="en",
            secret_word="xyznotinvocab123",  # Not in vocab
            top_n=10000
        )

        # Reset mock to ensure clean state
        mock_grpc_context.reset_mock()

        response = servicer.StartGame(request, mock_grpc_context)

        # Should set error code for word not in vocabulary
        mock_grpc_context.set_code.assert_called_once_with(grpc.StatusCode.INVALID_ARGUMENT)
        mock_grpc_context.set_details.assert_called_once()

    def test_start_game_normalizes_word(self, servicer, mock_grpc_context):
        """Test that secret word is normalized (lowercase, stripped)."""
        from proto.compute.v1 import service_pb2

        request = service_pb2.StartGameRequest(
            language="en",
            secret_word="  CAT  ",  # Mixed case with spaces
            top_n=10000
        )

        response = servicer.StartGame(request, mock_grpc_context)

        # Should succeed with normalized word
        assert response.calculation_time >= 0
        mock_grpc_context.set_code.assert_not_called()

    def test_start_game_default_top_n(self, servicer, mock_grpc_context):
        """Test default top_n when not specified or 0."""
        from proto.compute.v1 import service_pb2

        request = service_pb2.StartGameRequest(
            language="en",
            secret_word="cat",
            top_n=0  # Should default to 10000
        )

        response = servicer.StartGame(request, mock_grpc_context)

        # Should succeed
        assert response.calculation_time >= 0

    # --- MakeGuess Tests ---

    def test_make_guess_success(self, servicer, start_game_request, make_guess_request, mock_grpc_context):
        """Test successful guess with rank returned."""
        # First start a game to populate rankings
        servicer.StartGame(start_game_request, mock_grpc_context)

        # Now make a guess
        response = servicer.MakeGuess(make_guess_request, mock_grpc_context)

        # Should return rank (dog is at rank 2)
        assert response.rank == 2
        assert response.found is False
        assert response.in_vocabulary is True

    def test_make_guess_exact_match(self, servicer, start_game_request, mock_grpc_context):
        """Test guessing the exact word returns found=True."""
        from proto.compute.v1 import service_pb2

        # First start a game
        servicer.StartGame(start_game_request, mock_grpc_context)

        # Guess the exact word
        request = service_pb2.MakeGuessRequest(
            guess="cat",
            secret_word="cat",
            language="en"
        )

        response = servicer.MakeGuess(request, mock_grpc_context)

        assert response.rank == 1
        assert response.found is True
        assert response.in_vocabulary is True

    def test_make_guess_not_in_vocabulary(self, servicer, start_game_request, mock_grpc_context):
        """Test guess with unknown word returns in_vocabulary=False."""
        from proto.compute.v1 import service_pb2

        # First start a game to populate rankings
        servicer.StartGame(start_game_request, mock_grpc_context)

        # Clear the mock's call tracking for fresh assertions
        mock_grpc_context.reset_mock()

        request = service_pb2.MakeGuessRequest(
            guess="xyzunknownword",
            secret_word="cat",
            language="en"
        )

        response = servicer.MakeGuess(request, mock_grpc_context)

        # Word not in vocabulary should return in_vocabulary=False
        assert response.rank == 0
        assert response.found is False
        assert response.in_vocabulary is False

    def test_make_guess_normalizes_input(self, servicer, start_game_request, mock_grpc_context):
        """Test that guess is normalized (lowercase, stripped)."""
        from proto.compute.v1 import service_pb2

        # First start a game
        servicer.StartGame(start_game_request, mock_grpc_context)

        request = service_pb2.MakeGuessRequest(
            guess="  DOG  ",  # Mixed case with spaces
            secret_word="cat",
            language="en"
        )

        response = servicer.MakeGuess(request, mock_grpc_context)

        # Should still find dog at rank 2
        assert response.rank == 2
        assert response.in_vocabulary is True

    def test_make_guess_no_game_data(self, servicer, mock_grpc_context):
        """Test guess without starting game first."""
        from proto.compute.v1 import service_pb2

        # Don't start a game first
        request = service_pb2.MakeGuessRequest(
            guess="dog",
            secret_word="unknownsecret",
            language="en"
        )

        response = servicer.MakeGuess(request, mock_grpc_context)

        # Should set NOT_FOUND error
        mock_grpc_context.set_code.assert_called_once_with(grpc.StatusCode.NOT_FOUND)

    # --- HealthCheck Tests ---

    def test_health_check_healthy(self, servicer, mock_grpc_context):
        """Test health check returns OK with loaded languages."""
        from proto.compute.v1 import service_pb2

        request = service_pb2.HealthCheckRequest()

        response = servicer.HealthCheck(request, mock_grpc_context)

        assert response.status == "healthy"
        assert "en" in response.loaded_languages

    def test_health_check_multiple_languages(self, mock_word_model, mock_grpc_context):
        """Test health check shows all loaded languages."""
        from grpc_server.servicer import ComputeServiceServicer
        from proto.compute.v1 import service_pb2

        # Add another language to the mock
        mock_word_model.models["ua"] = MagicMock()

        servicer = ComputeServiceServicer(mock_word_model)
        request = service_pb2.HealthCheckRequest()

        response = servicer.HealthCheck(request, mock_grpc_context)

        assert response.status == "healthy"
        assert len(response.loaded_languages) == 2
        assert "en" in response.loaded_languages
        assert "ua" in response.loaded_languages
