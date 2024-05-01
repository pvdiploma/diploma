CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    secret TEXT NOT NULL
);

CREATE TABLE roles (
    role_id int PRIMARY KEY,
    description varchar(32) UNIQUE NOT NULL
);

INSERT INTO roles VALUES (0, 'admin'), (1, 'orginizer'), (2, 'distributor'), (3, 'buyer');


CREATE TABLE IF NOT EXISTS users (
    id       SERIAL PRIMARY KEY,
    login    TEXT NOT NULL,
    email    TEXT UNIQUE NOT NULL,
    pwd_hash bytea NOT NULL,
    role     INTEGER REFERENCES roles (role_id) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
