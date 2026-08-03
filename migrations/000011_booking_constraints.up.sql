ALTER TABLE seat_bookings
DROP CONSTRAINT IF EXISTS seat_bookings_show_id_seat_id_key;

ALTER TABLE seat_bookings
ADD CONSTRAINT seat_bookings_booking_id_seat_id_key
UNIQUE (booking_id, seat_id);