create table if not exists seats(
    id bigserial primary key,
    hall_id bigint not null,
    seat_number varchar(10) not null,
    seat_type varchar(20) default 'regular' check (seat_type in ('regular','standard','premium')),
    is_active boolean not null default true,
    created_at timestamptz default now() not null,
    updated_at timestamptz default now() not null,
    unique(hall_id,seat_number)
);