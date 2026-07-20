create table if not exists seat_bookings(
id bigserial primary key,
booking_id bigint not null references bookings(id),
seat_id bigint not null references seats(id),
show_id bigint not null references shows(id),
created_at  timestamptz not null default now(),
unique(show_id,seat_id)
);