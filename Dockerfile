FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./bin ./bin
COPY ./config ./config
COPY ./internals ./internals

RUN CGO_ENABLED=0 go build -o ./rightshift ./bin/main.go

# Stage 2
FROM alpine:latest

WORKDIR /app/

COPY --from=0 /app/rightshift ./
COPY ./migrations/ ./migrations/

EXPOSE 8000

CMD ["./rightshift"]