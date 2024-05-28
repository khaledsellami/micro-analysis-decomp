FROM golang:1.21-alpine3.20 AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/MicroAnalyzer

FROM alpine:3.20

COPY --from=builder /go/bin/MicroAnalyzer .

ENTRYPOINT [ "./MicroAnalyzer" ]