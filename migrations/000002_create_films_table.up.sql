create table if not exists films
(
    id          serial primary key,
    user_id     int references users (id),
    title       varchar not null,
    year        int,
    genre       varchar,
    description varchar,
    rating      float,
    photo_url   varchar,
    comment     varchar,
    is_viewed   boolean                     default false,
    user_rating float,
    review      varchar,
    created_at  timestamp(5) with time zone default NOW(),
    updated_at  timestamp(5) with time zone default NOW()
);