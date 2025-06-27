FROM golang:1.22-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/ws-ai-advice/cmd
RUN go build -o ws-ai-advice

CMD ["./ws-ai-advice"]