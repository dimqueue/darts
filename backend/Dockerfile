FROM golang:1.24-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS build-dev
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3 && \
    rm -rf docs/ && \
    swag init -g cmd/app/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -tags dev -o /usr/local/bin/darts github.com/dimqueue/darts/cmd/app

FROM base AS build-prod
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/darts github.com/dimqueue/darts/cmd/app

FROM alpine:3.19 AS dev
RUN apk add --no-cache ca-certificates
COPY --from=build-dev /usr/local/bin/darts /usr/local/bin/darts
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]

FROM alpine:3.19 AS prod
RUN apk add --no-cache ca-certificates
COPY --from=build-prod /usr/local/bin/darts /usr/local/bin/darts
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]