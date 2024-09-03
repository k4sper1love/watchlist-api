CREATE TABLE IF NOT EXISTS films
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT                NOT NULL,
    title       TEXT                     NOT NULL,
    year        INT,
    genre       TEXT,
    description TEXT,
    rating      FLOAT,
    photo_url   TEXT,
    comment     TEXT,
    is_viewed   BOOLEAN                  NOT NULL DEFAULT FALSE,
    user_rating FLOAT,
    review      TEXT,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);