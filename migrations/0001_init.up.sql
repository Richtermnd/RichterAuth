CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO roles (name) VALUES ('superadmin'), ('admin'), ('user');


CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY NOT NULL,
    username CHAR(32) NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    role_id INTEGER REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS confirm_keys (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    confirm_key TEXT NOT NULL,
    expired_at TIMESTAMP NOT NULL
);

--- Create superadmin user ---
--- Username: Admin ---
--- Email: admin@admin.com ---
--- Password: admin ---
--- !!! CREATE NEW SUPERADMIN OR REPLACE THIS BEFORE PRODUCTION!!! ---

INSERT INTO USERS (username, email, hashed_password, is_confirmed, role_id) 
    VALUES (
        'admin', 
        'admin@admin.com', 
        '$2a$10$HUMztg6bLIvcBr0bYB.Wl.uEBOuCJ/8KaV3wkNyCfTEa5H5S8hmre',
        TRUE,
        (SELECT id FROM roles WHERE name = 'superadmin')
        );
