# Dockerfile for user service

FROM golang:1.13

LABEL maintainer="Ryoei Shimaura"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8000

CMD ["./main"]