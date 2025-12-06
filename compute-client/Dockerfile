FROM python:3.11-slim

WORKDIR /app

# Install dependencies first (better caching)
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY src/ ./src/

# Create logs directory
RUN mkdir -p logs

# Unbuffered output for better logging
ENV PYTHONUNBUFFERED=1

# Expose both ports (only one will be used based on SERVER_MODE)
EXPOSE 5000 50051

# Run the application
CMD ["python", "-m", "src.main"]