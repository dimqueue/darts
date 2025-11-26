import grpc
import logging
from concurrent import futures

from proto.compute.v1 import service_pb2_grpc
from models.word_similarity import WordSimilarityModel
from .servicer import ComputeServiceServicer

logger = logging.getLogger(__name__)


def serve(word_model: WordSimilarityModel, port: int = 50051):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service_pb2_grpc.add_ComputeServiceServicer_to_server(
        ComputeServiceServicer(word_model), server
    )

    server.add_insecure_port(f'[::]:{port}')
    server.start()

    logger.info(f"gRPC server started on port {port}")
    return server


if __name__ == "__main__":
    from config.settings import Config
    from utils.logging_config import setup_logging
    from utils.monitor import ResourceMonitor

    setup_logging(service_name="compute-grpc", log_dir=Config.LOG_DIR)

    logger.info("=" * 60)
    logger.info("Starting gRPC server in standalone mode")
    logger.info("=" * 60)

    # Initialize monitoring
    monitor = ResourceMonitor()
    monitor.log_metrics(context="Startup")

    # Initialize model
    logger.info(f"Loading model: {Config.MODEL_NAME}")
    word_model = WordSimilarityModel()
    word_model.load_model(Config.DEFAULT_LANGUAGE, Config.MODEL_NAME)
    monitor.log_metrics(context="Model Loaded")

    logger.info("Warming up model...")
    word_model.warm_up(Config.DEFAULT_LANGUAGE)
    monitor.log_metrics(context="Model Warmed Up")

    # Start server
    server = serve(word_model, port=Config.GRPC_PORT)

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Shutting down gRPC server...")
        server.stop(0)
        word_model.cleanup()
        monitor.log_metrics(context="Shutdown")