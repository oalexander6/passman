FROM node:22-alpine AS clientbuild

WORKDIR /app

COPY ./web/package.json ./
COPY ./web/package-lock.json ./

RUN npm install --omit=dev

COPY ./web ./

RUN npm run build

FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o ./passman ./cmd/main.go

# Stage 2
FROM alpine:latest

WORKDIR /app/

COPY --from=1 /app/passman ./
COPY --from=0 /app/dist ./web/dist

EXPOSE 8000

CMD ["./passman"]