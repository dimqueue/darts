FROM golang:1.23-alpine AS builder

ARG BUILD_MODE=prod

WORKDIR /app

# Install swag only for dev mode
RUN if [ "$BUILD_MODE" = "dev" ]; then \
        go install github.com/swaggo/swag/cmd/swag@v1.16.3; \
    fi

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate swagger docs only for dev mode
RUN if [ "$BUILD_MODE" = "dev" ]; then \
        rm -rf docs/ && \
        swag init -g cmd/app/main.go; \
    fi

# Build with dev tags only for dev mode
RUN if [ "$BUILD_MODE" = "dev" ]; then \
        CGO_ENABLED=0 GOOS=linux go build -tags dev -o /usr/local/bin/darts github.com/dimqueue/darts/cmd/app; \
    else \
        CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/darts github.com/dimqueue/darts/cmd/app; \
    fi

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/darts /usr/local/bin/darts

EXPOSE 8080

ENTRYPOINT ["darts"]