
-- +goose Up 

ALTER TABLE users
ADD COLUMN access_token TEXT ;


