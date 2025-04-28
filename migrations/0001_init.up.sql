-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending'
);

-- +goose Down
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;
