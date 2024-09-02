create table if not exists collections
(
    id          serial primary key,
    user_id     int references users (id),
    name        varchar,
    description varchar,
    created_at  timestamp(5) with time zone default NOW(),
    updated_at  timestamp(5) with time zone default NOW()
);
