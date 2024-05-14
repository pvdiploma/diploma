CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER REFERENCES users (id) NOT NULL,
    name varchar(255) NOT NULL,
    description TEXT NOT NULL,
    country varchar(255) NOT NULL,
    city varchar(255) NOT NULL,
    place varchar(255) NOT NULL,
    address varchar(255) NOT NULL,
    date DATE NOT NULL,
    ticket_amount INTEGER NOT NULL,
    age varchar(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_event_owner_id ON events (owner_id);

CREATE TABLE IF NOT EXISTS event_categories (
    id SERIAL PRIMARY KEY,
    event_id INTEGER REFERENCES events (id) NOT NULL,
    category varchar(255) NOT NULL,
    price INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    UNIQUE (event_id, category)  
);

CREATE INDEX IF NOT EXISTS idx_event_category_event_id ON event_categories (event_id);