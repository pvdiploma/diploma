CREATE TYPE deal_status AS ENUM ('accepted', 'rejected', 'pending');

CREATE TABLE IF NOT EXISTS deals (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER REFERENCES users (id) NOT NULL,
    recipient_id INTEGER REFERENCES users (id) NOT NULL,
    organizer_id INTEGER REFERENCES users (id) NOT NULL,
    distributor_id INTEGER REFERENCES users (id) NOT NULL,
    event_id INTEGER REFERENCES events (id) NOT NULL,
    commission INTEGER NOT NULL,
    status deal_status NOT NULL,

    UNIQUE (organizer_id, distributor_id, event_id)
)

CREATE TABLE IF NOT EXISTS widgets (
    id SERIAL PRIMARY KEY,
    deal_id INTEGER REFERENCES deals (id) NOT NULL,
    body TEXT NOT NULL,
    script TEXT NOT NULL
)