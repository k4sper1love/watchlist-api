create table if not exists refresh_tokens (
    id serial primary key,
    token text not null,
    user_id int references users(id) on delete cascade ,
    expires_at timestamp with time zone,
    revoked boolean default false
);