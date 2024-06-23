# Building the binary of the App
FROM golang:1.22.4

COPY go.mod go.sum /

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Create the `public` dir and copy all the assets into it
RUN mkdir ./static
COPY ./static ./static

RUN apt update && apt install -y git
