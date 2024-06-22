
-- +goose Up 
-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT UNIQUE NOT NULL ,
    post_url TEXT UNIQUE NOT NULL ,
    content TEXT   NOT NULL,
    thumbnail TEXT 
   
);

-- +goose Down
DROP TABLE posts;


