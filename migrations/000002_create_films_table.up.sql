CREATE TYPE view_status AS ENUM ('not_viewed', 'in_progress', 'viewed');

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
    view_status view_status DEFAULT 'not_viewed',
    user_rating NUMERIC(4, 2),
    review      TEXT,
    url TEXT,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
