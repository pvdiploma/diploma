FROM golang:1.21-bullseye

RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

COPY go.mod .
COPY go.sum .

RUN go mod tidy && go mod verify

COPY . . 