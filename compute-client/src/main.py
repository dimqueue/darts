import logging
import sys

from config.settings import Config
from utils.logging_config import setup_logging
from models.word_similarity import WordSimilarityModel
from utils.monitor import ResourceMonitor

logger = logging.getLogger(__name__)


def initialize_model(monitor: ResourceMonitor) -> WordSimilarityModel:
    logger.info(f"Loading model: {Config.MODEL_NAME} for language: {Config.DEFAULT_LANGUAGE}")
    monitor.log_metrics(context="Before Model Load")

    word_model = WordSimilarityModel()
    word_model.load_model(Config.DEFAULT_LANGUAGE, Config.MODEL_NAME)
    monitor.log_metrics(context="Model Loaded")

    logger.info("Warming up model...")
    word_model.warm_up(Config.DEFAULT_LANGUAGE)
    monitor.log_metrics(context="Model Warmed Up")

    return word_model


def run_http_server(word_model: WordSimilarityModel, monitor: ResourceMonitor):
    import uvicorn
    from api.app import create_app

    app = create_app(word_model, monitor)

    try:
        uvicorn.run(app, host=Config.HTTP_HOST, port=Config.HTTP_PORT, log_config=None)
    except KeyboardInterrupt:
        logger.info("HTTP server shutting down...")
        word_model.cleanup()
        monitor.log_metrics(context="Shutdown")


def run_grpc_server(word_model: WordSimilarityModel, monitor: ResourceMonitor):
    import threading
    from grpc_server.server import serve as serve_grpc

    server = serve_grpc(word_model, port=Config.GRPC_PORT)

    # Start periodic monitoring in a background thread
    stop_monitoring = threading.Event()

    def periodic_monitor():
        logger.info(f"Periodic monitoring started (interval: {Config.MONITORING_INTERVAL}s)")
        while not stop_monitoring.wait(Config.MONITORING_INTERVAL):
            monitor.log_metrics(context="Periodic")
        logger.info("Periodic monitoring stopped")

    monitoring_thread = threading.Thread(target=periodic_monitor, daemon=True)
    monitoring_thread.start()

    logger.info("gRPC server ready")

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("gRPC server shutting down...")
        stop_monitoring.set()
        server.stop(0)
        word_model.cleanup()
        monitor.log_metrics(context="Shutdown")


def main():
    setup_logging(service_name="compute-client", log_dir=Config.LOG_DIR)

    mode = Config.SERVER_MODE.upper()
    port = Config.GRPC_PORT if Config.SERVER_MODE == "grpc" else Config.HTTP_PORT
    host = f"0.0.0.0:{port}" if Config.SERVER_MODE == "grpc" else f"{Config.HTTP_HOST}:{Config.HTTP_PORT}"

    logger.info("=" * 60)
    logger.info(f"Compute Client - {mode} Server")
    logger.info("=" * 60)
    logger.info(f"  Mode:     {mode}")
    logger.info(f"  Address:  {host}")
    logger.info(f"  Model:    {Config.MODEL_NAME}")
    logger.info(f"  Language: {Config.DEFAULT_LANGUAGE}")
    logger.info("=" * 60)

    monitor = ResourceMonitor()
    monitor.log_metrics(context="Startup")

    try:
        word_model = initialize_model(monitor)
    except Exception as e:
        logger.error(f"Failed to initialize model: {e}")
        sys.exit(1)

    if Config.SERVER_MODE == "http":
        run_http_server(word_model, monitor)
    elif Config.SERVER_MODE == "grpc":
        run_grpc_server(word_model, monitor)
    else:
        logger.error(f"Invalid SERVER_MODE: {Config.SERVER_MODE}")
        sys.exit(1)


if __name__ == "__main__":
    main()