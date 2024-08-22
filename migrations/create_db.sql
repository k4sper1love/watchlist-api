-- create database watchlist;
-- drop database if exists watchlist;
-- \c watchlist;

drop table if exists user_collections;
drop table if exists collection_films;
drop table if exists collections;
drop table if exists user_viewed;
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
    creator_id int default 1,
    title varchar not null,
    year int,
    genre varchar,
    description varchar,
    rating float,
    photo_url varchar,
    foreign key (creator_id) references users(id) on delete set default
);

create table if not exists user_viewed (
    user_id int not null,
    film_id int not null,
    rating float,
    review varchar,
    viewed_at timestamp default current_timestamp,
    primary key (user_id, film_id),
    foreign key (user_id) references users(id) on delete cascade ,
    foreign key (film_id) references films(id)
);

create table if not exists collections (
    id serial primary key,
    creator_id int references users(id),
    title varchar,
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

create table if not exists user_collections (
    user_id int not null,
    collection_id int not null,
    added_at timestamp default current_timestamp,
    primary key (user_id, collection_id),
    foreign key (user_id) references users(id),
    foreign key (collection_id) references collections(id)
)