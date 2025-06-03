FROM golang:1.22

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/user-service
RUN go build -o user ./cmd/main.go

EXPOSE 50051
CMD ["./user"]