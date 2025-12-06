FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY src/ ./src/

RUN mkdir -p logs

ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH=/app/src

EXPOSE 5000 50051

# Run the application
CMD ["python", "-m", "src.main"]