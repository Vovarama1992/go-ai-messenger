CREATE TABLE chat_members (
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    accepted BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (chat_id, user_id),
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
);