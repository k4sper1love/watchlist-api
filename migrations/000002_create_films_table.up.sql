CREATE TABLE IF NOT EXISTS films
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT                NOT NULL,
    is_favorite   BOOLEAN                  NOT NULL DEFAULT FALSE,
    title       TEXT                     NOT NULL,
    year        INT,
    genre       TEXT,
    description TEXT,
    rating      NUMERIC(4, 2),
    image_url   TEXT,
    comment     TEXT,
    is_viewed   BOOLEAN                  NOT NULL DEFAULT FALSE,
    user_rating NUMERIC(4, 2),
    review      TEXT,
    url TEXT,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);