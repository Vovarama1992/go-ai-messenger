Документация по микросервисам

1. Общие правила

📁 Расположение интерфейсов (ports)

Все зависимости между слоями (gRPC-клиенты, Kafka-паблишеры и пр.) описываются через интерфейсы и располагаются в:

internal/**/ports/ Это тербует команд make generate-mocks

Интерфейсы используются в usecase-слое, конкретные реализации — в adapters.

🌐 HTTP-роуты и Swagger

Все HTTP-роуты описываются в функции RegisterRoutes, например:

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
    // @Summary Логин
    // @Description Аутентификация пользователя и выдача JWT
    // ...
    mux.HandleFunc("/login", handler.Login)
}

Файл с роутами должен называться routes.go и располагаться в одной из папок:

internal/**/http

Swagger-аннотации размещаются над mux.HandleFunc(...) в этом файле. Это требует команда make swagger-init

2. auth-service

📱 HTTP API (наружу)

POST /login — аутентификация пользователя

POST /register — регистрация нового пользователя


🔌 gRPC API (внутрь системы)

rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);

📁 Определение: authpb/auth.proto🔹 Сервер регистрируется в auth-service

🔗 Зависимости

auth-service по gRPC обращается к user-service:

CreateUser(email, passwordHash) → userID

GetByEmail(email) → userID, passwordHash

