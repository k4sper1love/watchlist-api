create table if not exists permissions
(
    id   serial primary key,
    code varchar unique not null
);

create table if not exists user_permissions
(
    user_id        int not null,
    permissions_id int not null,
    primary key (user_id, permissions_id),
    foreign key (user_id) references users (id) on delete cascade,
    foreign key (permissions_id) references permissions (id)
);

insert into permissions (code)
values ('film:create'),
       ('collection:create');