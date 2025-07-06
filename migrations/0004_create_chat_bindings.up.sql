CREATE TABLE chat_bindings (
    chat_id BIGINT PRIMARY KEY,
    thread_id TEXT,
    user_id BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('advice', 'autoreply')),
    created_at BIGINT NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
);