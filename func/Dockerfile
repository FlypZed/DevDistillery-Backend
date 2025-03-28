FROM golang:1.21 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/main .
COPY --from=build /app/.env .
EXPOSE 8080
ENTRYPOINT ["./main"]