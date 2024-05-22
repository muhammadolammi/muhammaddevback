
-- +goose Up 

ALTER TABLE tutorials
ADD COLUMN thumbnail TEXT ;


-- +goose Down

ALTER TABLE tutorials
DROP COLUMN thumbnail  ;