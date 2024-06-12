-- +goose Up 
-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token TEXT UNIQUE NOT NULL ,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,

    user_id UUID REFERENCES users(id)  ON DELETE CASCADE NOT NULL

);



