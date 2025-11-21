import logging
import os
from datetime import datetime
from pathlib import Path
from logging.handlers import RotatingFileHandler


def setup_logging(
    service_name: str = "compute-client",
    log_dir: str = "logs",
    max_bytes: int = 10 * 1024 * 1024,  # 10 MB
    backup_count: int = 5
):
    """
    Configure logging to write to both console and rotating log files.
    Creates a new log file for each service start with timestamp.

    Args:
        service_name: Name of the service (used in log filename)
        log_dir: Directory where log files will be stored
        max_bytes: Maximum size per log file before rotation (default: 10MB)
        backup_count: Number of backup files to keep (default: 5)

    Returns:
        Path to the created log file
    """

    log_path = Path(log_dir)
    log_path.mkdir(exist_ok=True)

    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    log_filename = log_path / f"{service_name}_{timestamp}.log"

    formatter = logging.Formatter(
        fmt='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S'
    )

    root_logger = logging.getLogger()
    root_logger.setLevel(logging.INFO)

    root_logger.handlers.clear()

    console_handler = logging.StreamHandler()
    console_handler.setLevel(logging.INFO)
    console_handler.setFormatter(formatter)
    root_logger.addHandler(console_handler)

    file_handler = RotatingFileHandler(
        log_filename,
        maxBytes=max_bytes,
        backupCount=backup_count,
        encoding='utf-8'
    )
    file_handler.setLevel(logging.INFO)
    file_handler.setFormatter(formatter)
    root_logger.addHandler(file_handler)

    logging.getLogger('gensim').setLevel(logging.WARNING)
    logging.getLogger('urllib3').setLevel(logging.WARNING)
    logging.getLogger('asyncio').setLevel(logging.WARNING)

    logging.info(f"Logging initialized. Log file: {log_filename}")
    logging.info(f"Service: {service_name}")
    logging.info(f"Log rotation: {max_bytes / 1024 / 1024:.0f}MB per file, {backup_count} backups")

    return log_filename