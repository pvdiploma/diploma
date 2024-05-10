CREATE TABLE IF NOT EXISTS ticket (
    id SERIAL PRIMARY KEY,
    event_category_id INTEGER REFERENCES event_category (id) NOT NULL,
    name varchar(255) NOT NULL,
    surname varchar(255) NOT NULL,
    patronymic varchar(255),
    email varchar(255) NOT NULL,
    discount INTEGER NOT NULL,
    total INTEGER NOT NULL,
    qr_code bytea,
    is_activated BOOLEAN NOT NULL,
    image_bytes bytea NOT NULL,
    image_path varchar(255) NOT NULL
);