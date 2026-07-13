FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o booking-system ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/booking-system .


EXPOSE 8080
CMD [ "./booking-system" ]