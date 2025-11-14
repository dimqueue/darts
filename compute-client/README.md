# Compute Client

Python microservice for computational tasks using FastAPI and Gensim.

## Structure

```
compute-client/
├── src/              # Source code
├── configs/          # Configuration files
├── Dockerfile        # Container definition
├── requirements.txt  # Python dependencies for Docker
└── environment.yml   # Conda environment for local dev
```

## Local Development (with Conda)

```bash
# Create/update conda environment
conda env create -f environment.yml
# or update existing:
conda env update -f environment.yml

# Activate environment
conda activate darts-compute

# Run locally
cd src
uvicorn main:app --reload --host 0.0.0.0 --port 5000
```

## Docker Deployment

```bash
# Build
docker compose build compute-client

# Run
docker compose up compute-client
```

## API Endpoints

- `GET /health` - Health check
- `POST /compute` - Computation endpoint (add your logic)