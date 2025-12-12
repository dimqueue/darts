FROM python:3.11-slim

WORKDIR /app

RUN useradd -m -u 1000 appuser

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY src/ ./src/

RUN mkdir -p logs && chown -R appuser:appuser /app

ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH=/app/src

USER appuser

EXPOSE 5000 50051

# Run the application
CMD ["python", "-m", "src.main"]