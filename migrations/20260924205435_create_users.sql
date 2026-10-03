-- +goose Up

CREATE TABLE users (
   id BIGSERIAL PRIMARY KEY,
   email VARCHAR(200) NOT NULL UNIQUE,
   hashed_password VARCHAR(200) NOT NULL
);

-- +goose Down

DROP TABLE users;