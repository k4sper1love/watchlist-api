CREATE TABLE IF NOT EXISTS films
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT                NOT NULL,
    title       TEXT                     NOT NULL,
    year        INT,
    genre       TEXT,
    description TEXT,
    rating      FLOAT,
    image_url   TEXT,
    comment     TEXT,
    is_viewed   BOOLEAN                  NOT NULL DEFAULT FALSE,
    user_rating FLOAT,
    review      TEXT,
    url TEXT,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);