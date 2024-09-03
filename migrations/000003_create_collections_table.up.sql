CREATE TABLE IF NOT EXISTS collections
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT                   NOT NULL,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);