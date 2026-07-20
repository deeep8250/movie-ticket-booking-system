create table if not exists theaters(
    id bigserial primary key,
    theater_name varchar not null,
    theater_owner varchar not null,
    theater_email varchar not null unique,
    city varchar not null,
    pin_code varchar not null,
    state varchar not null,
    district varchar not null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now()  not null
);