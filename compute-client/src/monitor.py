import psutil
import logging
from typing import Dict, Optional

logger = logging.getLogger(__name__)


class ResourceMonitor:
    """
    Monitors system resource usage for the application.
    Designed to be protocol-agnostic (works with HTTP, gRPC, etc.)
    """

    def __init__(self):
        self.process = psutil.Process()
        self.initial_memory = self.process.memory_info().rss / 1024 / 1024
        self.process.cpu_percent(interval=None)
        logger.info(f"ResourceMonitor initialized. Initial memory: {self.initial_memory:.2f}MB")

    def get_metrics(self, include_cpu: bool = True) -> Dict[str, float]:
        """
        Get current resource usage metrics.

        Args:
            include_cpu: Whether to measure CPU (adds small delay, default: True)

        Returns:
            Dictionary with memory, CPU, threads, and file descriptors metrics
        """
        try:
            with self.process.oneshot():
                memory_info = self.process.memory_info()

                metrics = {
                    "memory_mb": memory_info.rss / 1024 / 1024,
                    "memory_percent": self.process.memory_percent(),
                    "num_threads": self.process.num_threads(),
                }

                if include_cpu:
                    cpu = self.process.cpu_percent(interval=None)
                    metrics["cpu_percent"] = cpu if cpu > 0 else 0.0
                else:
                    metrics["cpu_percent"] = 0.0

                # num_fds only available on Unix systems
                if hasattr(self.process, 'num_fds'):
                    try:
                        metrics["num_fds"] = self.process.num_fds()
                    except (psutil.AccessDenied, psutil.NoSuchProcess):
                        metrics["num_fds"] = 0

                return metrics
        except Exception as e:
            logger.warning(f"Failed to collect metrics: {e}")
            return {
                "memory_mb": 0.0,
                "memory_percent": 0.0,
                "cpu_percent": 0.0,
                "num_threads": 0,
            }

    def log_metrics(self, context: str = "Periodic", include_cpu: bool = True) -> Dict[str, float]:

        metrics = self.get_metrics(include_cpu=include_cpu)

        log_msg = (
            f"[{context}] "
            f"Memory: {metrics['memory_mb']:.1f}MB ({metrics['memory_percent']:.1f}%), "
            f"Threads: {metrics['num_threads']}"
        )

        if include_cpu and metrics['cpu_percent'] > 0:
            log_msg += f", CPU: {metrics['cpu_percent']:.1f}%"

        logger.info(log_msg)

        return metrics

    def get_memory_delta(self) -> float:
        """Get memory change since initialization in MB"""
        current_memory = self.process.memory_info().rss / 1024 / 1024
        return current_memory - self.initial_memory