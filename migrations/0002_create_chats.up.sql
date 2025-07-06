CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    creator_id BIGINT NOT NULL,
    chat_type VARCHAR(20) NOT NULL CHECK (type IN ('private', 'group')),
    created_at BIGINT NOT NULL
);