FROM golang:latest

WORKDIR /server
COPY . /server

RUN go mod download