create extension if not exists citext;

create table if not exists users (
    id serial primary key,
    username citext not null unique,
    email citext not null unique,
    password varchar not null,
    created_at timestamp(5) with time zone default NOW(),
    version int not null default 1
);