FROM golang:1.14.4-alpine AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build
RUN go build -o /memeater ./cmd/memeater/main.go
RUN go build -o /cpuburner ./cmd/cpuburner/main.go

FROM golang:1.14.4-alpine

COPY --from=builder /app/px-usewithcare /
COPY --from=builder /memeater ./bin
COPY --from=builder /cpuburner ./bin


ENTRYPOINT /px-usewithcare
