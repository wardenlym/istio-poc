FROM golang:alpine as builder

WORKDIR /app
COPY main.go .

RUN go build -o golang-app-gateway main.go

FROM alpine:latest
RUN apk add curl
WORKDIR /app
COPY --from=builder /app/golang-app-gateway .
COPY math.jfif math.jfif
CMD ["./golang-app-gateway"]