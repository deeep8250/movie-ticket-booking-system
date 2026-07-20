create table if not exists movies(
    id bigserial primary key,
    title varchar(200)  not null,
    description text not null,
    language varchar(200) not null,
    duration_min int not null check(duration_min>0),
    release_date date not null,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null
);