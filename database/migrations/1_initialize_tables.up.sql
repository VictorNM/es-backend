CREATE TABLE users
(
    id              int generated always as identity,
    email           varchar(255) unique,
    username        varchar(255) unique,
    hashed_password varchar(255) not null,
    full_name       varchar(255) not null,
    phone           varchar(13),
    year_of_birth   integer,
    is_active       boolean default false,
    is_super_admin  boolean default false,
    country         varchar(2),
    gender          varchar(6),
    language        varchar(2),
    activation_key  varchar(255),
    created_at      timestamp,
    updated_at      timestamp,

    primary key (id)
);