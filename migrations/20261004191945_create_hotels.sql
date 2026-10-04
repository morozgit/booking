-- +goose Up
CREATE TABLE hotels (
   id BIGSERIAL PRIMARY KEY,
   title VARCHAR(100),
   location VARCHAR(400)
);

-- +goose Down
DROP TABLE hotels;
