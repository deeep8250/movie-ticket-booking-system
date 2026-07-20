    create table if not exists bookings(
        id bigserial primary key,
        user_id bigint references users(id) not null,
        show_id bigint references shows(id) not null,
        status varchar(20) not null default 'pending' check(status in('confirmed','cancelled','pending','expired')),
        created_at timestamptz default now(),
        updated_at timestamptz default now()
    );