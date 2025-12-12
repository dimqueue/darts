import logging
import sys
from datetime import datetime
from pathlib import Path
from logging.handlers import RotatingFileHandler


class ColorCodes:

    RESET = "\033[0m"
    BOLD = "\033[1m"

    # Foreground colors
    BLACK = "\033[30m"
    RED = "\033[31m"
    GREEN = "\033[32m"
    YELLOW = "\033[33m"
    BLUE = "\033[34m"
    MAGENTA = "\033[35m"
    CYAN = "\033[36m"
    WHITE = "\033[37m"

    # Bright foreground colors
    BRIGHT_BLACK = "\033[90m"
    BRIGHT_RED = "\033[91m"
    BRIGHT_GREEN = "\033[92m"
    BRIGHT_YELLOW = "\033[93m"
    BRIGHT_BLUE = "\033[94m"
    BRIGHT_MAGENTA = "\033[95m"
    BRIGHT_CYAN = "\033[96m"
    BRIGHT_WHITE = "\033[97m"


class ColoredFormatter(logging.Formatter):

    LEVEL_COLORS = {
        'DEBUG': ColorCodes.BRIGHT_BLACK,
        'INFO': ColorCodes.BRIGHT_CYAN,
        'WARNING': ColorCodes.BRIGHT_YELLOW,
        'ERROR': ColorCodes.BRIGHT_RED,
        'CRITICAL': ColorCodes.RED + ColorCodes.BOLD,
    }

    def __init__(self, use_colors=True):
        super().__init__()
        self.use_colors = use_colors and self._supports_color()

    def _supports_color(self):
        if not hasattr(sys.stdout, 'isatty') or not sys.stdout.isatty():
            return False
        return True

    def format(self, record):

        if not hasattr(record, 'funcName'):
            record.funcName = 'unknown'
        if not hasattr(record, 'lineno'):
            record.lineno = 0

        timestamp = self.formatTime(record, '%Y-%m-%d %H:%M:%S')

        if self.use_colors:
            level_color = self.LEVEL_COLORS.get(record.levelname, '')
            reset = ColorCodes.RESET
            time_color = ColorCodes.BRIGHT_BLACK
            name_color = ColorCodes.BLUE
            location_color = ColorCodes.BRIGHT_BLACK
        else:
            level_color = reset = time_color = name_color = location_color = ''

        level = f"{record.levelname:<8}"

        name = record.name
        if len(name) > 25:
            name = '...' + name[-22:]

        file_line = f"{record.filename}:{record.lineno}"
        func_name = f"{record.funcName}()"

        if len(file_line) > 25:
            file_line = '...' + file_line[-22:]
        if len(func_name) > 25:
            func_name = func_name[:22] + '...()'

        location = f"{file_line:<25} in {func_name:<25}"

        parts = [
            f"{time_color}{timestamp}{reset}",
            f"{level_color}{level}{reset}",
            f"{name_color}{name:<25}{reset}",
            f"{location_color}{location}{reset}",
            record.getMessage()
        ]

        return " | ".join(parts)


class PlainFormatter(logging.Formatter):

    def format(self, record):
        timestamp = self.formatTime(record, '%Y-%m-%d %H:%M:%S')
        level = f"{record.levelname:<8}"
        name = record.name
        location = f"{record.filename}:{record.lineno}"
        func = record.funcName
        message = record.getMessage()

        return f"{timestamp} | {level} | {name:<25} | {location:<30} | {func:<20} | {message}"


def setup_logging(
    service_name: str = "compute-client",
    log_dir: str = "logs",
    log_level: str = "DEBUG",
    max_bytes: int = 10 * 1024 * 1024,  # 10 MB
    backup_count: int = 5,
    use_colors: bool = True,
    log_format: str = "colored"  # "colored" or "json" (for future)
):

    log_path = Path(log_dir)
    log_path.mkdir(exist_ok=True)

    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    log_filename = log_path / f"{service_name}_{timestamp}.log"

    level = getattr(logging, log_level.upper(), logging.INFO)

    # Create formatters
    console_formatter = ColoredFormatter(use_colors=use_colors)
    file_formatter = PlainFormatter()

    # Configure root logger
    root_logger = logging.getLogger()
    root_logger.setLevel(level)
    root_logger.handlers.clear()

    # Console handler (colored)
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setLevel(level)
    console_handler.setFormatter(console_formatter)
    root_logger.addHandler(console_handler)

    # File handler (plain structured)
    file_handler = RotatingFileHandler(
        log_filename,
        maxBytes=max_bytes,
        backupCount=backup_count,
        encoding='utf-8'
    )
    file_handler.setLevel(level)
    file_handler.setFormatter(file_formatter)
    root_logger.addHandler(file_handler)

    # Suppress noisy third-party loggers
    logging.getLogger('gensim').setLevel(logging.WARNING)
    logging.getLogger('urllib3').setLevel(logging.WARNING)
    logging.getLogger('asyncio').setLevel(logging.WARNING)
    logging.getLogger('uvicorn.access').setLevel(logging.WARNING)

    # Log initialization info
    logging.info(f"Logging initialized. File: {log_filename}")
    logging.info(f"Service: {service_name} | Colors: {use_colors}")

    return log_filename