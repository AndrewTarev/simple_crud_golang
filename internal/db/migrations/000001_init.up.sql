CREATE TABLE users
(
    id serial primary key,
    username varchar(255) not null unique,
    email varchar(255) unique,
    password varchar(255) not null,
    created_at timestamp default now()
);