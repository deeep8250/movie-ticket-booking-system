create table if not exists users(
 id bigserial primary key,
 username varchar not null,
 email varchar(200) not null unique,
 mobile int(20) not null unique,
 created_at timestamptz not null default now(),
 updated_at timestamptz not null default now() 
);
