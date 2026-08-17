BEGIN;

TRUNCATE TABLE
    seat_bookings,
    bookings,
    seats,
    shows,
    movies,
    halls,
    theaters,
    users
RESTART IDENTITY CASCADE;

INSERT INTO theaters (
    id,
    theater_name,
    theater_owner,
    theater_email,
    city,
    pin_code,
    state,
    district
)
VALUES (
    1,
    'Test Theater',
    'Test Owner',
    'test-theater@example.com',
    'Test City',
    '123456',
    'Test State',
    'Test District'
);

INSERT INTO halls (
    id,
    theater_id,
    hall_name
)
VALUES (
    1,
    1,
    'Test Hall 1'
);

INSERT INTO movies (
    id,
    title,
    description,
    language,
    duration_min,
    release_date
)
VALUES (
    1,
    'Test Movie',
    'Test movie for integration tests',
    'English',
    120,
    CURRENT_DATE
);

INSERT INTO shows (
    id,
    movie_id,
    hall_id,
    starts_at,
    ends_at,
    base_price,
    status
)
VALUES
(
    1,
    1,
    1,
    NOW() + INTERVAL '1 day',
    NOW() + INTERVAL '1 day 2 hours',
    450,
    'active'
),
(
    2,
    1,
    1,
    NOW() + INTERVAL '2 days',
    NOW() + INTERVAL '2 days 2 hours',
    450,
    'active'
);

INSERT INTO seats (
    id,
    hall_id,
    seat_number,
    seat_type,
    is_active
)
SELECT
    gs,
    1,
    'A' || gs,
    'regular',
    true
FROM generate_series(1, 20) AS gs;

SELECT setval(pg_get_serial_sequence('theaters', 'id'), COALESCE((SELECT MAX(id) FROM theaters), 1));
SELECT setval(pg_get_serial_sequence('halls', 'id'), COALESCE((SELECT MAX(id) FROM halls), 1));
SELECT setval(pg_get_serial_sequence('movies', 'id'), COALESCE((SELECT MAX(id) FROM movies), 1));
SELECT setval(pg_get_serial_sequence('shows', 'id'), COALESCE((SELECT MAX(id) FROM shows), 1));
SELECT setval(pg_get_serial_sequence('seats', 'id'), COALESCE((SELECT MAX(id) FROM seats), 1));

COMMIT;