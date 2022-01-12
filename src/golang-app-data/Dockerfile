FROM golang:alpine as builder

WORKDIR /app
COPY main.go .

RUN go build -o golang-app-data main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/golang-app-data .
CMD ["./golang-app-data"]