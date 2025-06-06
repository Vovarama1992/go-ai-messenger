FROM golang:1.22

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

WORKDIR /app/ai-service
RUN go build -o ai ./cmd/main.go

EXPOSE 8085

CMD ["./ai"]