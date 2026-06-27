FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6
COPY . .
RUN swag init && go build -o personas-api .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/personas-api .
EXPOSE 8080
CMD ["./personas-api"]
