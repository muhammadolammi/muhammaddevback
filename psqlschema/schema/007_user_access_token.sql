
-- +goose Up 

ALTER TABLE users
ADD COLUMN access_token TEXT ;


-- +goose Down

ALTER TABLE users
DROP COLUMN access_token  ;