FROM golang:1.23 AS builder

WORKDIR /go/src/github.com/dimqueue/darts

COPY . .

COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/darts github.com/dimqueue/darts/cmd/app

FROM alpine:3.19

COPY --from=builder /usr/local/bin/darts /usr/local/bin/darts
RUN apk add --no-cache ca-certificates

EXPOSE 8080

ENTRYPOINT ["darts"]
