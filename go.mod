module github.com/Vovarama1992/go-ai-messenger

go 1.23.0

require (
	github.com/go-chi/chi/v5 v5.2.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/golang/mock v1.6.0
	github.com/googollee/go-socket.io v1.7.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/lib/pq v1.10.9
	github.com/segmentio/kafka-go v0.4.48
	github.com/stretchr/testify v1.10.0
	go.uber.org/mock v0.5.2
	golang.org/x/crypto v0.37.0
	google.golang.org/grpc v1.54.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gomodule/redigo v1.8.4 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Vovarama1992/go-ai-messenger/proto/userpb => ./proto/userpb

replace github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/mock => ./message-service/internal/message/mock
