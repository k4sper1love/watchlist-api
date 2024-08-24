-- create database watchlist;
-- drop database if exists watchlist;
-- \c watchlist;
-- create extension if not exists "uuid-ossp";

drop table if exists collection_films;
drop table if exists collections;
drop table if exists viewed_films;
drop table if exists user_films;
drop table if exists films ;
drop table if exists users;

create table if not exists users (
    id serial primary key,
    username varchar not null unique,
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
    comment varchar,
    added_at timestamp default current_timestamp,
    primary key (collection_id, film_id),
    foreign key (collection_id) references collections(id),
    foreign key (film_id) references films(id)
);

create table if not exists viewed_films (
    user_id int not null,
    film_id int not null,
    rating float,
    review varchar,
    viewed_at timestamp default current_timestamp,
    primary key (user_id, film_id),
    foreign key (user_id) references users(id),
    foreign key (film_id) references films(id)
);

insert into users (username) values ('admin');
