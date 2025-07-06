ОБЩИЕ ПРАВИЛА

📁 Расположение интерфейсов (ports)

Все зависимости между слоями (gRPC-клиенты, Kafka-паблишеры и пр.) описываются через интерфейсы и располагаются в:

internal/**/ports/ Это требует команда make generate-mocks

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

ДОКУМЕНТАЦИЯ ПО МИКРОСЕРВИСАМ

1. auth-service

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



2. user-service

🔌 gRPC API

user-service предоставляет gRPC-интерфейс:

GetUserByEmail(email) → userID, passwordHash

CreateUser(email, passwordHash) → userID

3. chat-service
📱 HTTP API (наружу)

Все эндпоинты описаны через Swagger-аннотации в internal/**/http/routes.go.

Сгенерированная документация используется фронтом и внешними клиентами.

📌 Для генерации: make swagger-init

🔌 gRPC API

Интерфейс chatpb.ChatService предоставляет:

GetChatByID(chat_id) → Chat

GetBindingsByChat(chat_id) → List<Binding>

GetUserWithChatByThreadID(thread_id) → {userID, chatID, email}

GetUsersByChatID(chat_id) → List<userID>

GetThreadContext(thread_id) → {senderID, senderEmail, chatID, participants}

📁 Определение: proto/chatpb/chat.proto
🔹 Реализация сервера: internal/chat/adapters/grpc

🔗 Зависимости

gRPC-запрос к user-service: GetUserByID(userID)
gRPC-запрос к message-service: GetMessagesByChat(chatID)

Kafka-события:

TOPIC_AI_ADVICE_REQUEST — отправка запроса в AI

TOPIC_CHAT_INVITE — отправка инвайтов в вебсокет

gRPC-вызовы из ws-ai-advice:

получение участника по threadID

получение контекста чата по threadID

4. message-service
📡 Kafka Listener (внутрь)
Читает сообщения из Kafka-топика:

TOPIC_MESSAGE_PERSIST — сохранение входящих сообщений в БД

💬 gRPC API (наружу)

GetMessagesByChat(chatID) → List<Message>
Возвращает историю сообщений по указанному чату

🤖 Обработка AI-сообщений

Если сообщение приходит только с ThreadID, без ChatID и SenderID,
то оно считается сгенерированным AI. В этом случае:

SenderID и ChatID подставляются через GetUserWithChatByThreadID(threadID)

AIGenerated = true

🔗 Зависимости

gRPC-запросы к:

chat-service: GetUserWithChatByThreadID(threadID)

user-service: GetUserByID(userID)