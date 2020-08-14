FROM golang:1.14.4-alpine AS builder
RUN apk add build-base libstdc++ --no-cache
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build

ENTRYPOINT ./px-usewithcare
