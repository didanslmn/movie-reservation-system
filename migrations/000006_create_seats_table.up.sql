BEGIN;
CREATE TABLE seats (
    id SERIAL PRIMARY KEY,
    cinema_hall_id INT NOT NULL REFERENCES cinema_hall(id) ON DELETE CASCADE,
    seat_number VARCHAR(10) NOT NULL,
    row VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
COMMIT;