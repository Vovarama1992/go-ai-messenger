CREATE TABLE IF NOT EXISTS messages (
    id           BIGSERIAL PRIMARY KEY,
    chat_id      BIGINT NOT NULL REFERENCES chats(id),
    sender_id    BIGINT NOT NULL REFERENCES users(id),
    content      TEXT NOT NULL,
    ai_generated BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Индексы на chat_id и sender_id для фильтрации
CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages(chat_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);