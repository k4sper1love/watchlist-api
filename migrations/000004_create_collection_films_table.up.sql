create table if not exists collection_films
(
    collection_id int not null,
    film_id       int not null,
    created_at    timestamp(5) with time zone default NOW(),
    updated_at    timestamp(5) with time zone default NOW(),
    primary key (collection_id, film_id),
    foreign key (collection_id) references collections (id),
    foreign key (film_id) references films (id)
);