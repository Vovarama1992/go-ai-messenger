# go-ai-messenger

## 🧠 Цель проекта

Реализовать мессенджер с поддержкой AI-ответов (советы, автореакция) в микросервисной архитектуре, с масштабируемыми WebSocket-шлюзами и Kafka-пайплайном.
в котором:
- пользователь может создать чат;
- чат может быть привязан к AI-сервису (совет или автоответ);
- AI-сервис знает, какой чат к какому пользователю относится;
- AI-сервис отвечает либо советом, либо автоматически вместо пользователя.

---

⚙️ Стек технологий

Go 1.22

gRPC

Kafka (Confluent)

PostgreSQL

WebSocket (socket.io-style)

Docker Compose (dev) / Kubernetes (prod)

Swagger (через HTTP шлюзы)

Prometheus + Grafana (мониторинг и логирование)

⚙️ Продакшн и масштабирование

Stateless-сервисы: легко реплицируются

Kafka — для async обмена сообщениями и AI

gRPC — для sync-взаимодействия сервисов

WebSocket-шлюз (ws-gateway) — масштабируется с sticky session или Redis pub/sub

Docker Compose — для локальной разработки

Kubernetes — целевая среда

CI/мейк: make test, make lint, make run, make migrate

Swagger подключается ко всем HTTP-сервисам

Прометей + Графана: сбор метрик и наблюдаемость

---

## 🧱 Микросервисы (gRPC, Kafka, PostgreSQL)

| Сервис          | Ответственность                    | gRPC-клиенты               | Kafka                            |
| --------------- | ---------------------------------- | -------------------------- | -------------------------------- |
| auth-service    | JWT-аутентификация, login/register | user-service               | ❌                                |
| user-service    | Хранение пользователей             | —                          | ❌                                |
| chat-service    | Чаты, AI-привязка                  | auth-service, user-service | ❌                                |
| message-service | Сохранение сообщений               | chat-service (optional)    | ✅ consume `chat.message.persist` |

### WebSocket

| Сервис       | Kafka topics                              | Примечание                                                           |
| ------------ | ----------------------------------------- | -------------------------------------------------------------------- |
| ws-gateway   | produce: persist, forward consume: forward | Один сервис, масштабируется (replicas) + sticky session/Redis pubsub |
| ws-ai-advice | consume: ai-advice                        | Пуш ответа AI-совета в `user:{id}`                                   |
| ws-autoreply | consume: forward                          | Фильтр AI-автоответов (senderId == targetUserId)                     |

### AI Service

| Задача                           | Kafka topics                                                                |
| -------------------------------- | --------------------------------------------------------------------------- |
| AI-обработка (советы/автоответы) | consume: ai.advice-request, ai.autoreply-request  produce: ai-advice, forward |

---

## 🔄 Kafka Topics

| Топик                             | Producer               | Consumer                 |
| --------------------------------- | ---------------------- | ------------------------ |
| chat.message.persist              | ws-gateway             | message-service          |
| chat.message.forward              | ws-gateway, ai-service | ws-gateway, ws-autoreply |
| chat.message.ai-advice            | ai-service             | ws-ai-advice             |
| chat.message.ai.advice-request    | ws-gateway             | ai-service               |
| chat.message.ai.autoreply-request | ws-gateway             | ai-service               |

---

## 📅 .env.example

```env
# Postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=go_messenger
POSTGRES_PORT=5432
POSTGRES_HOST=postgres
DATABASE_URL=postgres://postgres:postgres@postgres:5432/go_messenger?sslmode=disable

# gRPC ports
USER_GRPC_PORT=50051
AUTH_HTTP_PORT=8080
AUTH_GRPC_PORT=50052
CHAT_HTTP_PORT=8081
CHAT_GRPC_PORT=50053

# Kafka
KAFKA_BROKER=kafka:9092

# Worker scaling
CHAT_MESSAGE_PERSIST_WORKER_COUNT=4
```

---

## 📚 Структура репозитория

```
go-ai-messenger/
├── ai-service/
├── auth-service/
├── chat-service/
├── message-service/
│   ├── cmd/
│   └── internal/
│       ├── delivery/kafka/
│       ├── usecase/
│       ├── infra/{kafka,postgres}
│       ├── ports/
│       └── model/
├── user-service/
├── ws-gateway/              # (в планах)
├── ws-ai-advice/            # (в планах)
├── ws-autoreply/            # (в планах)
├── proto/
├── docker-compose.yml
├── .env.example
├── go.mod
├── MANIFEST.md
└── README.md
```

---

## 📊 Уже готово

| Сервис          | Юнит-тесты           | Интеграционные | Комментарий                           |
| --------------- | -------------------- | -------------- | ------------------------------------- |
| auth-service    | ✅ с моком userClient | в планах       |                                       |
| user-service    | ✅ с мок repo         | в планах       |                                       |
| chat-service    | в разработке         | —              |                                       |
| message-service | пока нет             | пока нет       | запланирован Kafka consumer + persist |

---

## ⚙️ Продакшн и масштабирование

- Stateless-сервисы: легко реплицируются
- Kafka — для async обмена сообщениями и AI
- gRPC — для sync-взаимодействия сервисов
- WebSocket-шлюз (`ws-gateway`) — масштабируется с sticky session или Redis pub/sub
- Docker Compose — для локальной разработки
- Kubernetes — целевая среда
- CI/мейк: `make test`, `make lint`, `make run`, `make migrate`

НЕ ЗАБЫТЬ В КОНЦЕ ОБЕРНУТЬ ВСЕ В СВАГГЕР И НАПИСТАЬ МИГРАЦИИ И МЕЙК ФАЙЛ 
