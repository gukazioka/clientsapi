FROM golang:1.23.0-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/main.go /app

COPY data.csv /app

CMD ["./main.go"]