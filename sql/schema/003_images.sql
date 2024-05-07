-- +goose Up 
-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    image_url TEXT UNIQUE NOT NULL
    );

    
-- +goose Down
DROP TABLE images;
