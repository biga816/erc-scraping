FROM golang:latest

ARG SRC_PATH="src/erc-scraping"

WORKDIR $GOPATH/$SRC_PATH
ADD . $GOPATH/$SRC_PATH

RUN go mod download