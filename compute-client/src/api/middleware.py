import logging
import time
import uuid

from fastapi import Request

logger = logging.getLogger(__name__)


async def log_request_resources(request: Request, call_next):
    request_id = str(uuid.uuid4())[:8]
    request.state.request_id = request_id
    start_time = time.time()

    # Log incoming request details
    client_host = request.client.host if request.client else "unknown"
    query_params = dict(request.query_params) if request.query_params else {}

    logger.info(
        f"[{request_id}] Incoming request: {request.method} {request.url.path} | "
        f"Client: {client_host} | "
        f"Query params: {query_params if query_params else 'none'}"
    )

    try:
        response = await call_next(request)
        duration = time.time() - start_time
        metrics = request.app.state.monitor.get_metrics(include_cpu=False)

        logger.info(
            f"[{request_id}] {request.method} {request.url.path} | "
            f"Status: {response.status_code} | "
            f"Duration: {duration:.3f}s | "
            f"Memory: {metrics['memory_mb']:.1f}MB"
        )

        response.headers["X-Request-ID"] = request_id
        return response

    except Exception as e:
        duration = time.time() - start_time
        logger.error(
            f"[{request_id}] {request.method} {request.url.path} | "
            f"ERROR: {str(e)} | Duration: {duration:.3f}s"
        )
        raise
