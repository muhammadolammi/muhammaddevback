-- +goose Up 
-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE tutorials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT UNIQUE NOT NULL ,
    tutorial_url TEXT UNIQUE NOT NULL ,
    description TEXT   NOT NULL,
    youtube_link TEXT  NOT NULL,
    playlist_id UUID REFERENCES playlists(id)  ON DELETE CASCADE

);

-- +goose Down
DROP TABLE tutorials;

