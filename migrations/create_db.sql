-- create database watchlist;
-- drop database if exists watchlist;
-- \c watchlist;
-- create extension if not exists "uuid-ossp";

drop table if exists refresh_tokens;
drop table if exists collection_films;
drop table if exists collections;
drop table if exists viewed_films;
drop table if exists user_films;
drop table if exists films ;
drop table if exists users;

create table if not exists users (
    id serial primary key,
    username varchar not null unique,
    email varchar not null unique,
    password varchar not null,
    created_at timestamp default current_timestamp
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
    created_at timestamp default current_timestamp
);

create table if not exists collections (
    id serial primary key,
    user_id int references users(id),
    name varchar,
    description varchar,
    created_at timestamp default current_timestamp
);


create table if not exists collection_films (
    collection_id int not null,
    film_id int not null,
    added_at timestamp default current_timestamp,
    primary key (collection_id, film_id),
    foreign key (collection_id) references collections(id),
    foreign key (film_id) references films(id)
);

create table if not exists refresh_tokens (
    id serial primary key,
    token text not null,
    user_id int references users(id),
    expires_at timestamp,
    revoked boolean default false
);
