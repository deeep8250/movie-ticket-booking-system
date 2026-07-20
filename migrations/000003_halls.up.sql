create table if not exists halls(
    id bigserial primary key,
    theater_id bigint references theaters(id) not null,
    hall_name varchar not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    unique(theater_id,hall_name)
);