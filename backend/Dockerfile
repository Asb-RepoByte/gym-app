FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o gym-app ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gym-app .
EXPOSE 8000
CMD ["./gym-app"]
