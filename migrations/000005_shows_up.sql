create table if not exists shows (
    id bigserial primary key,
    movie_id bigint not null references movies(id),
    hall_id bigint not null references halls(id),
    starts_at timestamptz not null ,
    ends_at timestamptz not null,
    base_price numeric(10,2) not null check(base_price>0),
    status varchar(20) not null default 'scheduled' check (status in('scheduled','cancelled','completed')),
    created_at timestamptz not null  default now(),
    updated_at timestamptz not null default now(),
   unique (hall_id,starts_at) ,
    check (ends_at > starts_at)
);