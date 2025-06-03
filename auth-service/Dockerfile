FROM golang:1.22

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/auth-service
RUN go build -o auth ./cmd/main.go

EXPOSE 8080
CMD ["./auth"]