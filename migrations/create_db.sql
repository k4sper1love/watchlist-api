-- create database watchlist;
-- drop database if exists watchlist;
-- \c watchlist;
-- create extension if not exists "uuid-ossp";
drop table if exists user_permissions;
drop table if exists permissions;
drop table if exists refresh_tokens;
drop table if exists collection_films;
drop table if exists collections;
drop table if exists viewed_films;
drop table if exists user_films;
drop table if exists films ;
drop table if exists users;

create extension if not exists citext;

create table if not exists users (
    id serial primary key,
    username citext not null unique,
    email citext not null unique,
    password varchar not null,
    created_at timestamp(5) with time zone default NOW(),
    version int not null default 1
);

create table if not exists films (
    id serial primary key,
    user_id int references users(id),
    title varchar not null,
    year int,
    genre varchar,
    description varchar,
    rating float,
    photo_url varchar,
    comment varchar,
    is_viewed boolean default false,
    user_rating float,
    review varchar,
    created_at timestamp(5) with time zone default NOW(),
    updated_at timestamp(5) with time zone default NOW()
);

create table if not exists collections (
    id serial primary key,
    user_id int references users(id),
    name varchar,
    description varchar,
    created_at timestamp(5) with time zone default NOW(),
    updated_at timestamp(5) with time zone default NOW()
);


create table if not exists collection_films (
    collection_id int not null,
    film_id int not null,
    created_at timestamp(5) with time zone default NOW(),
    updated_at timestamp(5) with time zone default NOW(),
    primary key (collection_id, film_id),
    foreign key (collection_id) references collections(id),
    foreign key (film_id) references films(id)
);

create table if not exists refresh_tokens (
    id serial primary key,
    token text not null,
    user_id int references users(id) on delete cascade ,
    expires_at timestamp with time zone,
    revoked boolean default false
);

create table if not exists permissions (
    id serial primary key,
    code varchar unique not null
);

create table if not exists user_permissions (
    user_id int not null,
    permissions_id int not null,
    primary key (user_id, permissions_id),
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (permissions_id) references permissions(id)
);

insert into permissions (code) values ('film:create'), ('collection:create');