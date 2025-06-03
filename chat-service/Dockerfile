FROM golang:1.22

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

WORKDIR /app/chat-service
RUN go build -o chat ./cmd/main.go

EXPOSE 8081
EXPOSE 50053

CMD ["./chat"]