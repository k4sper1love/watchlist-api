create database wishlist;
drop database if exists wishlist;
\c wishlist;

drop table if exists users;
create table if not exists users (
    id serial primary key,
    username varchar not null unique,
    created_at timestamp default current_timestamp
);

insert into users (username) values ('admin');


drop table if exists films;
create table if not exists films (
    id serial primary key,
    title varchar not null,
    year int,
    genre varchar,
    description varchar,
    rating float,
    photo_url varchar
);

drop table if exists wishlist;
create table if not exists wishlist (
    user_id int not null,
    film_id int not null,
    added_at timestamp default current_timestamp,
    primary key (user_id, film_id),
    foreign key (user_id) references users(id),
    foreign key (film_id) references films(id)
)
